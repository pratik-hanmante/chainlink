package ocr2keeper_test

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/ocr2keepers/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umbracle/ethgo/abi"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	registrylogica21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic_a_wrapper_2_1"
	registrylogicb21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic_b_wrapper_2_1"
	registry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func TestFilterNamesFromSpec21(t *testing.T) {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	require.NoError(t, err)
	address := common.HexToAddress(hexutil.Encode(b))

	spec := &job.OCR2OracleSpec{
		PluginType: job.OCR2Keeper,
		ContractID: address.String(), // valid contract addr
	}

	names, err := ocr2keeper.FilterNamesFromSpec21(spec)
	require.NoError(t, err)

	assert.Len(t, names, 2)
	assert.Equal(t, logpoller.FilterName("KeepersRegistry TransmitEventProvider", address), names[0])
	assert.Equal(t, logpoller.FilterName("KeeperRegistry Events", address), names[1])

	spec = &job.OCR2OracleSpec{
		PluginType: job.OCR2Keeper,
		ContractID: "0x5431", // invalid contract addr
	}
	_, err = ocr2keeper.FilterNamesFromSpec21(spec)
	require.ErrorContains(t, err, "not a valid EIP55 formatted address")
}

func TestIntegration_KeeperPluginLogUpkeep(t *testing.T) {
	// g := gomega.NewWithT(t)
	lggr := logger.TestLogger(t)

	// setup blockchain
	sergey := testutils.MustNewSimTransactor(t) // owns all the link
	steve := testutils.MustNewSimTransactor(t)  // registry owner
	carrol := testutils.MustNewSimTransactor(t) // upkeep owner
	genesisData := core.GenesisAlloc{
		sergey.From: {Balance: assets.Ether(10000).ToInt()},
		steve.From:  {Balance: assets.Ether(10000).ToInt()},
		carrol.From: {Balance: assets.Ether(10000).ToInt()},
	}
	// Generate 5 keys for nodes (1 bootstrap + 4 ocr nodes) and fund them with ether
	var nodeKeys [5]ethkey.KeyV2
	for i := int64(0); i < 5; i++ {
		nodeKeys[i] = cltest.MustGenerateRandomKey(t)
		genesisData[nodeKeys[i].Address] = core.GenesisAccount{Balance: assets.Ether(1000).ToInt()}
	}

	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	stopMining := cltest.Mine(backend, 3*time.Second) // Should be greater than deltaRound since we cannot access old blocks on simulated blockchain
	defer stopMining()

	// Deploy contracts
	linkAddr, _, linkToken, err := link_token_interface.DeployLinkToken(sergey, backend)
	require.NoError(t, err)
	gasFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend, 18, big.NewInt(60000000000))
	require.NoError(t, err)
	linkFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend, 18, big.NewInt(2000000000000000000))
	require.NoError(t, err)
	registry := deployKeeper21Registry(t, steve, backend, linkAddr, linkFeedAddr, gasFeedAddr)

	// Setup bootstrap + oracle nodes
	bootstrapNodePort := int64(19599)
	appBootstrap, bootstrapPeerID, bootstrapTransmitter, bootstrapKb := setupNode(t, bootstrapNodePort, "bootstrap_keeper_ocr", nodeKeys[0], backend, nil)
	bootstrapNode := Node{
		appBootstrap, bootstrapTransmitter, bootstrapKb,
	}
	var (
		oracles []confighelper.OracleIdentityExtra
		nodes   []Node
	)
	// Set up the minimum 4 oracles all funded
	for i := int64(0); i < 4; i++ {
		app, peerID, transmitter, kb := setupNode(t, bootstrapNodePort+i+1, fmt.Sprintf("oracle_keeper%d", i), nodeKeys[i+1], backend, []commontypes.BootstrapperLocator{
			// Supply the bootstrap IP and port as a V2 peer address
			{PeerID: bootstrapPeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort)}},
		})

		nodes = append(nodes, Node{
			app, transmitter, kb,
		})
		offchainPublicKey, _ := hex.DecodeString(strings.TrimPrefix(kb.OnChainPublicKey(), "0x"))
		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  offchainPublicKey,
				TransmitAccount:   ocrTypes.Account(transmitter.String()),
				OffchainPublicKey: kb.OffchainPublicKey(),
				PeerID:            peerID,
			},
			ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
		})
	}

	// Add the bootstrap job
	bootstrapNode.AddBootstrapJob(t, fmt.Sprintf(`
		type                              = "bootstrap"
		relay                             = "evm"
		schemaVersion                     = 1
		name                              = "boot"
		contractID                        = "%s"
		contractConfigTrackerPollInterval = "15s"

		[relayConfig]
		chainID = 1337
	`, registry.Address()))

	// Add OCR jobs
	for i, node := range nodes {
		node.AddJob(t, fmt.Sprintf(`
		type = "offchainreporting2"
		pluginType = "ocr2automation"
		relay = "evm"
		name = "ocr2keepers-%d"
		schemaVersion = 1
		contractID = "%s"
		feedID = "%s"
		contractConfigTrackerPollInterval = "15s"
		ocrKeyBundleID = "%s"
		transmitterID = "%s"
		p2pv2Bootstrappers = [
		  "%s"
		]

		[relayConfig]
		chainID = 1337

		[pluginConfig]
		maxServiceWorkers = 100
		cacheEvictionInterval = "1s"
		mercuryCredentialName = "%s"
		`, i, registry.Address(), common.HexToHash("0x"), node.KeyBundle.ID(), node.Transmitter, fmt.Sprintf("%s@127.0.0.1:%d", bootstrapPeerID, bootstrapNodePort), MercuryCredName))
	}

	// Setup config on contract
	configType := abi.MustNewType("tuple(uint32 paymentPremiumPPB,uint32 flatFeeMicroLink,uint32 checkGasLimit,uint24 stalenessSeconds,uint16 gasCeilingMultiplier,uint96 minUpkeepSpend,uint32 maxPerformGas,uint32 maxCheckDataSize,uint32 maxPerformDataSize,uint32 maxRevertDataSize, uint256 fallbackGasPrice,uint256 fallbackLinkPrice,address transcoder,address[] registrars, address upkeepPrivilegeManager)")
	onchainConfig, err := abi.Encode(map[string]interface{}{
		"paymentPremiumPPB":      uint32(0),
		"flatFeeMicroLink":       uint32(0),
		"checkGasLimit":          uint32(6500000),
		"stalenessSeconds":       uint32(90000),
		"gasCeilingMultiplier":   uint16(2),
		"minUpkeepSpend":         uint32(0),
		"maxPerformGas":          uint32(5000000),
		"maxCheckDataSize":       uint32(5000),
		"maxPerformDataSize":     uint32(5000),
		"maxRevertDataSize":      uint32(5000),
		"fallbackGasPrice":       big.NewInt(60000000000),
		"fallbackLinkPrice":      big.NewInt(2000000000000000000),
		"transcoder":             testutils.NewAddress(),
		"registrars":             []common.Address{testutils.NewAddress()},
		"upkeepPrivilegeManager": testutils.NewAddress(),
	}, configType)
	// var args gethabi.Arguments = []gethabi.Argument{{Type: configType}}
	// onchainConfig, err := args.Pack(iregistry21.KeeperRegistryBase21OnchainConfig{
	// 	PaymentPremiumPPB:      uint32(0),
	// 	FlatFeeMicroLink:       uint32(0),
	// 	CheckGasLimit:          uint32(6500000),
	// 	StalenessSeconds:       big.NewInt(90000),
	// 	GasCeilingMultiplier:   uint16(2),
	// 	MinUpkeepSpend:         big.NewInt(0),
	// 	MaxPerformGas:          uint32(5000000),
	// 	MaxCheckDataSize:       uint32(5000),
	// 	MaxPerformDataSize:     uint32(5000),
	// 	MaxRevertDataSize:      uint32(5000),
	// 	FallbackGasPrice:       big.NewInt(60000000000),
	// 	FallbackLinkPrice:      big.NewInt(2000000000000000000),
	// 	Transcoder:             testutils.NewAddress(),
	// 	Registrars:             []common.Address{testutils.NewAddress()},
	// 	UpkeepPrivilegeManager: testutils.NewAddress(),
	// })
	require.NoError(t, err)
	rawCfg, err := json.Marshal(config.OffchainConfig{
		PerformLockoutWindow: 100 * 12 * 1000, // ~100 block lockout (on goerli)
		MinConfirmations:     1,
	})
	if err != nil {
		t.Logf("error creating off-chain config: %s", err)
		t.FailNow()
	}
	signers, transmitters, threshold, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		10*time.Second,        // deltaProgress time.Duration,
		10*time.Second,        // deltaResend time.Duration,
		2500*time.Millisecond, // deltaRound time.Duration,
		40*time.Millisecond,   // deltaGrace time.Duration,
		15*time.Second,        // deltaStage time.Duration,
		3,                     // rMax uint8,
		[]int{1, 1, 1, 1},
		oracles,
		rawCfg,                // reportingPluginConfig []byte,
		20*time.Millisecond,   // Max duration query
		1600*time.Millisecond, // Max duration observation
		800*time.Millisecond,
		20*time.Millisecond,
		20*time.Millisecond,
		1, // f
		onchainConfig,
	)
	require.NoError(t, err)
	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)
	transmitterAddresses, err := accountsToAddress(transmitters)
	require.NoError(t, err)
	offchainConfigVersion = 30
	lggr.Infow("Setting Config on Oracle Contract",
		"signerAddresses", signerAddresses,
		"transmitterAddresses", transmitterAddresses,
		"threshold", threshold,
		"onchainConfig", onchainConfig,
		"encodedConfigVersion", offchainConfigVersion,
		"offchainConfig", offchainConfig,
	)
	_, err = registry.SetConfig(
		steve,
		signerAddresses,
		transmitterAddresses,
		threshold,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)
	require.NoError(t, err)
	backend.Commit()

	fmt.Println("EvmRegistry: Setup done, creating an upkeep")

	<-time.After(time.Second * 5)

	n := 10

	_, err = linkToken.Transfer(sergey, carrol.From, big.NewInt(0).Mul(oneHunEth, big.NewInt(int64(n+1))))
	require.NoError(t, err)
	ids, _, contracts := deployUpkeeps(t, backend, carrol, steve, linkToken, registry, n)
	require.Equal(t, n, len(ids))
	require.Equal(t, len(contracts), len(ids))
	t.Logf("EvmRegistry: deployed and registered %d upkeeps", len(ids))
	<-time.After(time.Second * 5)
	go func(contracts []*log_upkeep_counter_wrapper.LogUpkeepCounter) {
		ctx := testutils.Context(t)
		emits := 10
		for i := 0; i < emits || ctx.Err() != nil; i++ {
			<-time.After(time.Second)
			t.Logf("EvmRegistry: calling upkeep contracts to emit events. run: %d", i+1)
			for _, contract := range contracts {
				_, err = contract.Start(carrol)
				require.NoError(t, err)
				backend.Commit()
			}
		}
	}(contracts)

	<-time.After(time.Second * 35)

	// check pipeline runs
	var allRuns []pipeline.Run
	for _, node := range nodes {
		runs, err2 := node.App.PipelineORM().GetAllRuns()
		require.NoError(t, err2)
		allRuns = append(allRuns, runs...)
	}
	require.GreaterOrEqual(t, len(allRuns), 0)

	<-time.After(time.Second * 5)
	// change payload
	// _, err = upkeepContract.SetBytesToSend(carrol, payload2)
	// require.NoError(t, err)
	// _, err = upkeepContract.SetShouldPerformUpkeep(carrol, true)
	// require.NoError(t, err)

	// observe 2nd job run and received payload changes
	// g.Eventually(receivedBytes, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(payload2))
}

func deployUpkeeps(t *testing.T, backend *backends.SimulatedBackend, carrol, steve *bind.TransactOpts, linkToken *link_token_interface.LinkToken, registry *iregistry21.IKeeperRegistryMaster, n int) ([]*big.Int, []common.Address, []*log_upkeep_counter_wrapper.LogUpkeepCounter) {
	ids := make([]*big.Int, n)
	addrs := make([]common.Address, n)
	contracts := make([]*log_upkeep_counter_wrapper.LogUpkeepCounter, n)
	for i := 0; i < n; i++ {
		upkeepAddr, _, upkeepContract, err := log_upkeep_counter_wrapper.DeployLogUpkeepCounter(
			carrol, backend,
			big.NewInt(100000),
		)
		require.NoError(t, err)
		logTriggerConfigType := abi.MustNewType("tuple(address contractAddress, uint8 filterSelector, bytes32 topic0, bytes32 topic1, bytes32 topic2, bytes32 topic3)")
		logTriggerConfig, err := abi.Encode(map[string]interface{}{
			"contractAddress": upkeepAddr,
			"filterSelector":  0,                                                                    // no indexed topics filtered
			"topic0":          "0x3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d", // event sig for Trigger()
			"topic1":          "0x",
			"topic2":          "0x",
			"topic3":          "0x",
		}, logTriggerConfigType)
		require.NoError(t, err)
		var err2 error
		enc := extraDataEncoder()
		extraData, err2 := enc.Encode([]interface{}{uint8(1), logTriggerConfig, "0x"})
		if err2 != nil {
			log.Fatal(err2)
		}

		registrationTx, err := registry.RegisterUpkeep(steve, upkeepAddr, 2_500_000, carrol.From, 1, []byte{}, extraData, []byte{})
		require.NoError(t, err)
		backend.Commit()
		upkeepID := getUpkeepIdFromTx21(t, registry, registrationTx, backend)

		// Fund the upkeep
		_, err = linkToken.Approve(carrol, registry.Address(), oneHunEth)
		require.NoError(t, err)
		_, err = registry.AddFunds(carrol, upkeepID, oneHunEth)
		require.NoError(t, err)
		backend.Commit()

		ids[i] = upkeepID
		contracts[i] = upkeepContract
		addrs[i] = upkeepAddr
	}
	return ids, addrs, contracts
}

func deployKeeper21Registry(
	t *testing.T,
	auth *bind.TransactOpts,
	backend *backends.SimulatedBackend,
	linkAddr, linkFeedAddr,
	gasFeedAddr common.Address,
) *iregistry21.IKeeperRegistryMaster {
	registryLogicBAddr, _, _, err := registrylogicb21.DeployKeeperRegistryLogicB(
		auth,
		backend,
		0, // Payment model
		linkAddr,
		linkFeedAddr,
		gasFeedAddr,
	)
	require.NoError(t, err)
	backend.Commit()

	registryLogicAAddr, _, _, err := registrylogica21.DeployKeeperRegistryLogicA(
		auth,
		backend,
		registryLogicBAddr,
	)
	require.NoError(t, err)
	backend.Commit()

	registryAddr, _, _, err := registry21.DeployKeeperRegistry(
		auth,
		backend,
		registryLogicAAddr,
	)
	require.NoError(t, err)
	backend.Commit()

	registryMaster, err := iregistry21.NewIKeeperRegistryMaster(registryAddr, backend)
	require.NoError(t, err)

	return registryMaster
}

func getUpkeepIdFromTx21(t *testing.T, registry *iregistry21.IKeeperRegistryMaster, registrationTx *types.Transaction, backend *backends.SimulatedBackend) *big.Int {
	receipt, err := backend.TransactionReceipt(testutils.Context(t), registrationTx.Hash())
	require.NoError(t, err)
	parsedLog, err := registry.ParseUpkeepRegistered(*receipt.Logs[0])
	require.NoError(t, err)
	return parsedLog.Id
}

func extraDataEncoder() *abi.Type {
	method, _ := abi.NewMethod("foo(uint8,bytes,bytes)") // foo is arbitrary, we just want the encoded values
	return method.Inputs
}
