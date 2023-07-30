package functions

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/functions/generated/functions_coordinator"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/functions_router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type logPollerWrapper struct {
	utils.StartStopOnce

	routerContract      *functions_router.FunctionsRouter
	pluginConfig        config.PluginConfig
	client              client.Client
	logPoller           logpoller.LogPoller
	subscribers         map[string]evmRelayTypes.RouteUpdateSubscriber
	activeCoordinator   common.Address
	proposedCoordinator common.Address
	nextBlockRequests   int64
	nextBlockResponses  int64
	mu                  sync.RWMutex
	closeWait           sync.WaitGroup
	stopCh              utils.StopChan
	lggr                logger.Logger
}

var _ evmRelayTypes.LogPollerWrapper = &logPollerWrapper{}

func NewLogPollerWrapper(routerContractAddress common.Address, pluginConfig config.PluginConfig, client client.Client, logPoller logpoller.LogPoller, lggr logger.Logger) (evmRelayTypes.LogPollerWrapper, error) {
	routerContract, err := functions_router.NewFunctionsRouter(routerContractAddress, client)
	if err != nil {
		return nil, err
	}

	return &logPollerWrapper{
		routerContract: routerContract,
		pluginConfig:   pluginConfig,
		logPoller:      logPoller,
		client:         client,
		subscribers:    make(map[string]evmRelayTypes.RouteUpdateSubscriber),
		stopCh:         make(utils.StopChan),
		lggr:           lggr,
	}, nil
}

func (l *logPollerWrapper) Start(context.Context) error {
	return l.StartOnce("LogPollerWrapper", func() error {
		l.lggr.Info("starting LogPollerWrapper")
		if l.pluginConfig.ContractVersion == 0 {
			l.mu.Lock()
			defer l.mu.Unlock()
			l.activeCoordinator = l.routerContract.Address()
			l.proposedCoordinator = l.routerContract.Address()
		} else if l.pluginConfig.ContractVersion == 1 {
			// TODO l.logPoller.LatestBlock()
			l.closeWait.Add(1)
			go l.checkForRouteUpdates()
		}
		return nil
	})
}

func (l *logPollerWrapper) Close() error {
	return l.StopOnce("LogPollerWrapper", func() (err error) {
		l.lggr.Info("closing LogPollerWrapper")
		close(l.stopCh)
		l.closeWait.Wait()
		return nil
	})
}

func (l *logPollerWrapper) HealthReport() map[string]error {
	return make(map[string]error)
}

func (l *logPollerWrapper) Name() string {
	return "LogPollerWrapper"
}

func (l *logPollerWrapper) Ready() error {
	return nil
}

// methods of LogPollerWrapper
func (l *logPollerWrapper) LatestRequests() ([]evmRelayTypes.OracleRequest, error) {
	// TODO all coordinators
	latest, err := l.logPoller.LatestBlock()
	if err != nil {
		return nil, err
	}
	// TODO checks, mutexes, second coordinator, update latest
	logs, err := l.logPoller.Logs(l.nextBlockRequests, latest, functions_coordinator.FunctionsCoordinatorOracleRequest{}.Topic(), l.activeCoordinator)
	if err != nil {
		return nil, err
	}

	coordinatorContract, err := functions_coordinator.NewFunctionsCoordinator(l.activeCoordinator, l.client)
	if err != nil {
		return nil, err
	}

	results := []evmRelayTypes.OracleRequest{}
	for _, log := range logs {
		gethLog := log.ToGethLog()
		oracleRequest, err := coordinatorContract.ParseOracleRequest(gethLog)
		if err != nil {
			continue
		}
		results = append(results, evmRelayTypes.OracleRequest{
			RequestId:          oracleRequest.RequestId,
			RequestingContract: oracleRequest.RequestingContract,
			RequestInitiator:   oracleRequest.RequestInitiator,
			SubscriptionId:     oracleRequest.SubscriptionId,
			SubscriptionOwner:  oracleRequest.SubscriptionOwner,
			Data:               oracleRequest.Data,
			DataVersion:        oracleRequest.DataVersion,
			Flags:              oracleRequest.Flags,
			CallbackGasLimit:   oracleRequest.CallbackGasLimit,
		})
	}

	return results, nil
}

func (l *logPollerWrapper) LatestResponses() ([]evmRelayTypes.OracleResponse, error) {
	// TODO poll both active and proposed coordinators for responses and parse them
	return nil, nil
}

// "internal" method called only by EVM relayer components
func (l *logPollerWrapper) SubscribeToUpdates(subscriberName string, subscriber evmRelayTypes.RouteUpdateSubscriber) {
	if l.pluginConfig.ContractVersion == 0 {
		// in V0, immediately set contract address to Oracle contract and never update again
		err := subscriber.UpdateRoutes(l.routerContract.Address(), l.routerContract.Address())
		l.lggr.Errorw("LogPollerWrapper: Failed to update routes", "subscriberName", subscriberName, "error", err)
	} else if l.pluginConfig.ContractVersion == 1 {
		l.mu.Lock()
		defer l.mu.Unlock()
		l.subscribers[subscriberName] = subscriber

		// TODO remove when periodic updates are ready and integration tests are properly migrated
		_ = subscriber.UpdateRoutes(l.routerContract.Address(), l.routerContract.Address())
	}
}

func (l *logPollerWrapper) checkForRouteUpdates() {
	defer l.closeWait.Done()
	freqSec := l.pluginConfig.ContractUpdateCheckFrequencySec
	if freqSec == 0 {
		l.lggr.Errorw("ContractUpdateCheckFrequencySec is zero - route update checks disabled")
		return
	}

	updateOnce := func() {
		// NOTE: timeout == frequency here, could be changed to a separate config value
		timeoutCtx, cancel := utils.ContextFromChanWithTimeout(l.stopCh, time.Duration(l.pluginConfig.ContractUpdateCheckFrequencySec)*time.Second)
		defer cancel()
		active, proposed, err := l.getCurrentCoordinators(timeoutCtx)
		if err != nil {
			l.lggr.Errorw("error calling UpdateFromContract", "err", err)
		}
		l.handleRouteUpdate(active, proposed)
	}

	updateOnce() // update once right away
	ticker := time.NewTicker(time.Duration(freqSec) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-l.stopCh:
			return
		case <-ticker.C:
			updateOnce()
		}
	}
}

func (l *logPollerWrapper) getCurrentCoordinators(ctx context.Context) (common.Address, common.Address, error) {
	if l.pluginConfig.ContractVersion == 0 {
		return l.routerContract.Address(), l.routerContract.Address(), nil
	}
	var donId [32]byte
	copy(donId[:], []byte(l.pluginConfig.DONId))

	activeCoordinator, err := l.routerContract.GetContractById(&bind.CallOpts{
		Pending: false,
		Context: ctx,
	}, donId, false)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	proposedCoordinator, err := l.routerContract.GetContractById(&bind.CallOpts{
		Pending: false,
		Context: ctx,
	}, donId, true)
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	return activeCoordinator, proposedCoordinator, nil
}

func (l *logPollerWrapper) handleRouteUpdate(activeCoordinator common.Address, proposedCoordinator common.Address) {
	changed := false
	l.mu.Lock()
	if activeCoordinator != l.activeCoordinator || proposedCoordinator != l.proposedCoordinator {
		l.activeCoordinator = activeCoordinator
		l.proposedCoordinator = proposedCoordinator
		changed = true
	}
	l.mu.Unlock()
	if !changed {
		return
	}

	l.registerFilters(activeCoordinator)
	l.registerFilters(proposedCoordinator)

	for _, subscriber := range l.subscribers {
		err := subscriber.UpdateRoutes(activeCoordinator, proposedCoordinator)
		if err != nil {
			l.lggr.Errorw("LogPollerWrapper: Failed to update routes", "error", err)
		}
	}
}

func filterName(addr common.Address) string {
	return logpoller.FilterName("FunctionsLogPollerWrapper", addr.String())
}

func (l *logPollerWrapper) registerFilters(coordinatorAddress common.Address) error {
	return l.logPoller.RegisterFilter(
		logpoller.Filter{
			Name: filterName(coordinatorAddress),
			EventSigs: []common.Hash{
				functions_coordinator.FunctionsCoordinatorOracleRequest{}.Topic(),
				functions_coordinator.FunctionsCoordinatorOracleResponse{}.Topic(),
			},
			Addresses: []common.Address{coordinatorAddress},
		})
}
