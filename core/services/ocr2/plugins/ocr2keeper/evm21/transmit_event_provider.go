package evm

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	ocr2keepers "github.com/smartcontractkit/ocr2keepers/pkg"
	"github.com/smartcontractkit/ocr2keepers/pkg/encoding"
	pluginutils "github.com/smartcontractkit/ocr2keepers/pkg/util"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_log_automation"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type TransmitUnpacker interface {
	UnpackTransmitTxInput([]byte) ([]ocr2keepers.UpkeepResult, error)
}

type TransmitEventProvider struct {
	sync              utils.StartStopOnce
	mu                sync.RWMutex
	runState          int
	runError          error
	logger            logger.Logger
	logPoller         logpoller.LogPoller
	registryAddress   common.Address
	lookbackBlocks    int64
	registry          *iregistry21.IKeeperRegistryMaster
	client            evmclient.Client
	packer            TransmitUnpacker
	txCheckBlockCache *pluginutils.Cache[string]
	cacheCleaner      *pluginutils.IntervalCacheCleaner[string]
}

func TransmitEventProviderFilterName(addr common.Address) string {
	return logpoller.FilterName("KeepersRegistry TransmitEventProvider", addr)
}

func NewTransmitEventProvider(
	logger logger.Logger,
	logPoller logpoller.LogPoller,
	registryAddress common.Address,
	client evmclient.Client,
	lookbackBlocks int64,
) (*TransmitEventProvider, error) {
	var err error

	contract, err := iregistry21.NewIKeeperRegistryMaster(common.HexToAddress("0x"), client)
	if err != nil {
		return nil, err
	}

	keeperABI, err := abi.JSON(strings.NewReader(iregistry21.IKeeperRegistryMasterABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}
	logDataABI, err := abi.JSON(strings.NewReader(i_log_automation.ILogAutomationABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrABINotParsable, err)
	}

	// Add log filters for the log poller so that it can poll and find the logs that
	// we need.
	err = logPoller.RegisterFilter(logpoller.Filter{
		Name: TransmitEventProviderFilterName(contract.Address()),
		EventSigs: []common.Hash{
			iregistry21.IKeeperRegistryMasterUpkeepPerformed{}.Topic(),
			iregistry21.IKeeperRegistryMasterReorgedUpkeepReport{}.Topic(),
			iregistry21.IKeeperRegistryMasterInsufficientFundsUpkeepReport{}.Topic(),
			iregistry21.IKeeperRegistryMasterStaleUpkeepReport{}.Topic(),
		},
		Addresses: []common.Address{registryAddress},
	})
	if err != nil {
		return nil, err
	}

	return &TransmitEventProvider{
		logger:            logger,
		logPoller:         logPoller,
		registryAddress:   registryAddress,
		lookbackBlocks:    lookbackBlocks,
		registry:          contract,
		client:            client,
		packer:            NewEvmRegistryPackerV2_1(keeperABI, logDataABI),
		txCheckBlockCache: pluginutils.NewCache[string](time.Hour),
		cacheCleaner:      pluginutils.NewIntervalCacheCleaner[string](time.Minute),
	}, nil
}

func (c *TransmitEventProvider) Name() string {
	return c.logger.Name()
}

func (c *TransmitEventProvider) Start(ctx context.Context) error {
	return c.sync.StartOnce("AutomationTransmitEventProvider", func() error {
		c.mu.Lock()
		defer c.mu.Unlock()

		go c.cacheCleaner.Run(c.txCheckBlockCache)
		c.runState = 1
		return nil
	})
}

func (c *TransmitEventProvider) Close() error {
	return c.sync.StopOnce("AutomationRegistry", func() error {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.cacheCleaner.Stop()
		c.runState = 0
		c.runError = nil
		return nil
	})
}

func (c *TransmitEventProvider) Ready() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.runState == 1 {
		return nil
	}
	return c.sync.Ready()
}

func (c *TransmitEventProvider) HealthReport() map[string]error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.runState > 1 {
		c.sync.SvcErrBuffer.Append(fmt.Errorf("failed run state: %w", c.runError))
	}
	return map[string]error{c.Name(): c.sync.Healthy()}
}

func (c *TransmitEventProvider) Events(ctx context.Context) ([]ocr2keepers.TransmitEvent, error) {
	end, err := c.logPoller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get latest block from log poller", err)
	}

	// always check the last lookback number of blocks and rebroadcast
	// this allows the plugin to make decisions based on event confirmations
	logs, err := c.logPoller.LogsWithSigs(
		end-c.lookbackBlocks,
		end,
		[]common.Hash{
			iregistry21.IKeeperRegistryMasterUpkeepPerformed{}.Topic(),
		},
		c.registryAddress,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	}

	performed, err := c.unmarshalPerformLogs(logs)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal logs", err)
	}

	return c.performedToTransmitEvents(performed, end)
}

func (c *TransmitEventProvider) PerformLogs(ctx context.Context) ([]ocr2keepers.PerformLog, error) {
	end, err := c.logPoller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get latest block from log poller", err)
	}

	// always check the last lookback number of blocks and rebroadcast
	// this allows the plugin to make decisions based on event confirmations
	logs, err := c.logPoller.LogsWithSigs(
		end-c.lookbackBlocks,
		end,
		[]common.Hash{
			iregistry21.IKeeperRegistryMasterUpkeepPerformed{}.Topic(),
		},
		c.registryAddress,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	}

	performed, err := c.unmarshalPerformLogs(logs)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal logs", err)
	}

	vals := []ocr2keepers.PerformLog{}
	for _, p := range performed {
		upkeepId := ocr2keepers.UpkeepIdentifier(p.Id.String())
		checkBlockNumber, err := c.getCheckBlockNumberFromTxHash(p.TxHash, upkeepId)
		if err != nil {
			c.logger.Error("error while fetching checkBlockNumber from reorged report log: %w", err)
			continue
		}
		l := ocr2keepers.PerformLog{
			Key:             encoding.BasicEncoder{}.MakeUpkeepKey(checkBlockNumber, upkeepId),
			TransmitBlock:   BlockKeyHelper[int64]{}.MakeBlockKey(p.BlockNumber),
			TransactionHash: p.TxHash.Hex(),
			Confirmations:   end - p.BlockNumber,
		}
		vals = append(vals, l)
	}

	return vals, nil
}

func (c *TransmitEventProvider) StaleReportLogs(ctx context.Context) ([]ocr2keepers.StaleReportLog, error) {
	end, err := c.logPoller.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get latest block from log poller", err)
	}

	// always check the last lookback number of blocks and rebroadcast
	// this allows the plugin to make decisions based on event confirmations

	// ReorgedUpkeepReportLogs
	logs, err := c.logPoller.LogsWithSigs(
		end-c.lookbackBlocks,
		end,
		[]common.Hash{
			iregistry21.IKeeperRegistryMasterReorgedUpkeepReport{}.Topic(),
		},
		c.registryAddress,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	}
	reorged, err := c.unmarshalReorgUpkeepLogs(logs)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal reorg logs", err)
	}

	// StaleUpkeepReportLogs
	logs, err = c.logPoller.LogsWithSigs(
		end-c.lookbackBlocks,
		end,
		[]common.Hash{
			iregistry21.IKeeperRegistryMasterStaleUpkeepReport{}.Topic(),
		},
		c.registryAddress,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	}
	staleUpkeep, err := c.unmarshalStaleUpkeepLogs(logs)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal stale upkeep logs", err)
	}

	// InsufficientFundsUpkeepReportLogs
	logs, err = c.logPoller.LogsWithSigs(
		end-c.lookbackBlocks,
		end,
		[]common.Hash{
			iregistry21.IKeeperRegistryMasterInsufficientFundsUpkeepReport{}.Topic(),
		},
		c.registryAddress,
		pg.WithParentCtx(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to collect logs from log poller", err)
	}
	insufficientFunds, err := c.unmarshalInsufficientFundsUpkeepLogs(logs)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal insufficient fund upkeep logs", err)
	}

	vals := []ocr2keepers.StaleReportLog{}
	for _, r := range reorged {
		upkeepId := ocr2keepers.UpkeepIdentifier(r.Id.String())
		checkBlockNumber, err := c.getCheckBlockNumberFromTxHash(r.TxHash, upkeepId)
		if err != nil {
			c.logger.Error("error while fetching checkBlockNumber from reorged report log: %w", err)
			continue
		}
		l := ocr2keepers.StaleReportLog{
			Key:             encoding.BasicEncoder{}.MakeUpkeepKey(checkBlockNumber, upkeepId),
			TransmitBlock:   BlockKeyHelper[int64]{}.MakeBlockKey(r.BlockNumber),
			TransactionHash: r.TxHash.Hex(),
			Confirmations:   end - r.BlockNumber,
		}
		vals = append(vals, l)
	}
	for _, r := range staleUpkeep {
		upkeepId := ocr2keepers.UpkeepIdentifier(r.Id.String())
		checkBlockNumber, err := c.getCheckBlockNumberFromTxHash(r.TxHash, upkeepId)
		if err != nil {
			c.logger.Error("error while fetching checkBlockNumber from stale report log: %w", err)
			continue
		}
		l := ocr2keepers.StaleReportLog{
			Key:             encoding.BasicEncoder{}.MakeUpkeepKey(checkBlockNumber, upkeepId),
			TransmitBlock:   BlockKeyHelper[int64]{}.MakeBlockKey(r.BlockNumber),
			TransactionHash: r.TxHash.Hex(),
			Confirmations:   end - r.BlockNumber,
		}
		vals = append(vals, l)
	}
	for _, r := range insufficientFunds {
		upkeepId := ocr2keepers.UpkeepIdentifier(r.Id.String())
		checkBlockNumber, err := c.getCheckBlockNumberFromTxHash(r.TxHash, upkeepId)
		if err != nil {
			c.logger.Error("error while fetching checkBlockNumber from insufficient funds report log: %w", err)
			continue
		}
		l := ocr2keepers.StaleReportLog{
			Key:             encoding.BasicEncoder{}.MakeUpkeepKey(checkBlockNumber, upkeepId),
			TransmitBlock:   BlockKeyHelper[int64]{}.MakeBlockKey(r.BlockNumber),
			TransactionHash: r.TxHash.Hex(),
			Confirmations:   end - r.BlockNumber,
		}
		vals = append(vals, l)
	}

	return vals, nil
}

func (c *TransmitEventProvider) unmarshalPerformLogs(logs []logpoller.Log) ([]performed, error) {
	results := []performed{}

	for _, log := range logs {
		rawLog := log.ToGethLog()
		abilog, err := c.registry.ParseLog(rawLog)
		if err != nil {
			return results, err
		}

		switch l := abilog.(type) {
		case *iregistry21.IKeeperRegistryMasterUpkeepPerformed:
			if l == nil {
				continue
			}

			r := performed{
				Log:                                  log,
				IKeeperRegistryMasterUpkeepPerformed: *l,
			}

			results = append(results, r)
		}
	}

	return results, nil
}

func (c *TransmitEventProvider) performedToTransmitEvents(performed []performed, latestBlock int64) ([]ocr2keepers.TransmitEvent, error) {
	var err error
	vals := []ocr2keepers.TransmitEvent{}

	for _, p := range performed {
		var checkBlockNumber ocr2keepers.BlockKey
		upkeepId := ocr2keepers.UpkeepIdentifier(p.Id.Bytes())
		switch getUpkeepType(upkeepId) {
		case conditionTrigger:
			checkBlockNumber, err = c.getCheckBlockNumberFromTxHash(p.TxHash, upkeepId)
			if err != nil {
				c.logger.Error("error while fetching checkBlockNumber from perform report log: %w", err)
				continue
			}
		default:
		}
		vals = append(vals, ocr2keepers.TransmitEvent{
			Type:            ocr2keepers.PerformEvent,
			TransmitBlock:   BlockKeyHelper[int64]{}.MakeBlockKey(p.BlockNumber),
			Confirmations:   latestBlock - p.BlockNumber,
			TransactionHash: p.TxHash.Hex(),
			ID:              hex.EncodeToString(p.TriggerID[:]),
			UpkeepID:        upkeepId,
			CheckBlock:      checkBlockNumber,
		})
	}

	return vals, nil
}

func (c *TransmitEventProvider) unmarshalReorgUpkeepLogs(logs []logpoller.Log) ([]reorged, error) {
	results := []reorged{}

	for _, log := range logs {
		rawLog := log.ToGethLog()
		abilog, err := c.registry.ParseLog(rawLog)
		if err != nil {
			return results, err
		}

		switch l := abilog.(type) {
		case *iregistry21.IKeeperRegistryMasterReorgedUpkeepReport:
			if l == nil {
				continue
			}

			r := reorged{
				Log:                                      log,
				IKeeperRegistryMasterReorgedUpkeepReport: *l,
			}

			results = append(results, r)
		}
	}

	return results, nil
}

func (c *TransmitEventProvider) unmarshalStaleUpkeepLogs(logs []logpoller.Log) ([]staleUpkeep, error) {
	results := []staleUpkeep{}

	for _, log := range logs {
		rawLog := log.ToGethLog()
		abilog, err := c.registry.ParseLog(rawLog)
		if err != nil {
			return results, err
		}

		switch l := abilog.(type) {
		case *iregistry21.IKeeperRegistryMasterStaleUpkeepReport:
			if l == nil {
				continue
			}

			r := staleUpkeep{
				Log:                                    log,
				IKeeperRegistryMasterStaleUpkeepReport: *l,
			}

			results = append(results, r)
		}
	}

	return results, nil
}

func (c *TransmitEventProvider) unmarshalInsufficientFundsUpkeepLogs(logs []logpoller.Log) ([]insufficientFunds, error) {
	results := []insufficientFunds{}

	for _, log := range logs {
		rawLog := log.ToGethLog()
		abilog, err := c.registry.ParseLog(rawLog)
		if err != nil {
			return results, err
		}

		switch l := abilog.(type) {
		case *iregistry21.IKeeperRegistryMasterInsufficientFundsUpkeepReport:
			if l == nil {
				continue
			}

			r := insufficientFunds{
				Log: log,
				IKeeperRegistryMasterInsufficientFundsUpkeepReport: *l,
			}

			results = append(results, r)
		}
	}

	return results, nil
}

// Fetches the checkBlockNumber for a particular transaction and an upkeep ID. Requires a RPC call to get txData
// so this function should not be used heavily
func (c *TransmitEventProvider) getCheckBlockNumberFromTxHash(txHash common.Hash, id ocr2keepers.UpkeepIdentifier) (bk ocr2keepers.BlockKey, e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("recovered from panic in getCheckBlockNumberForUpkeep: %v", r)
		}
	}()

	// Check if value already exists in cache for txHash, id pair
	cacheKey := txHash.String() + "|" + string(id)
	if val, ok := c.txCheckBlockCache.Get(cacheKey); ok {
		return ocr2keepers.BlockKey(val), nil
	}

	var tx gethtypes.Transaction
	err := c.client.CallContext(context.Background(), &tx, "eth_getTransactionByHash", txHash)
	if err != nil {
		return "", err
	}

	txData := tx.Data()
	if len(txData) < 4 {
		return "", fmt.Errorf("error in getCheckBlockNumberForUpkeep, got invalid tx data %s", txData)
	}

	decodedReport, err := c.packer.UnpackTransmitTxInput(txData[4:]) // Remove first 4 bytes of function signature
	if err != nil {
		return "", err
	}

	for _, upkeep := range decodedReport {
		// TODO: the log provider should be in the evm package for isolation
		res, ok := upkeep.(EVMAutomationUpkeepResult21)
		if !ok {
			return "", fmt.Errorf("unexpected type")
		}

		if res.ID.String() == string(id) {
			bl := fmt.Sprintf("%d", res.Block)

			c.txCheckBlockCache.Set(cacheKey, bl, pluginutils.DefaultCacheExpiration)

			return ocr2keepers.BlockKey(bl), nil
		}
	}

	return "", fmt.Errorf("upkeep %s not found in tx hash %s", id, txHash)
}

type performed struct {
	logpoller.Log
	iregistry21.IKeeperRegistryMasterUpkeepPerformed
}

type reorged struct {
	logpoller.Log
	iregistry21.IKeeperRegistryMasterReorgedUpkeepReport
}

type staleUpkeep struct {
	logpoller.Log
	iregistry21.IKeeperRegistryMasterStaleUpkeepReport
}

type insufficientFunds struct {
	logpoller.Log
	iregistry21.IKeeperRegistryMasterInsufficientFundsUpkeepReport
}
