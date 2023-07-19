package evm

import (
	"database/sql"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	testoffchainaggregator2 "github.com/smartcontractkit/libocr/gethwrappers2/testocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmClientMocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestConfigPoller(t *testing.T) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)
	b := backends.NewSimulatedBackend(core.GenesisAlloc{
		user.From: {Balance: big.NewInt(1000000000000000000)}},
		5*ethconfig.Defaults.Miner.GasCeil)
	linkTokenAddress, _, _, err := link_token_interface.DeployLinkToken(user, b)
	require.NoError(t, err)
	accessAddress, _, _, err := testoffchainaggregator2.DeploySimpleWriteAccessController(user, b)
	require.NoError(t, err, "failed to deploy test access controller contract")
	ocrAddress, _, ocrContract, err := ocr2aggregator.DeployOCR2Aggregator(
		user,
		b,
		linkTokenAddress,
		big.NewInt(0),
		big.NewInt(10),
		accessAddress,
		accessAddress,
		9,
		"TEST",
		false,
	)
	require.NoError(t, err)
	b.Commit()

	db := pgtest.NewSqlxDB(t)
	cfg := pgtest.NewQConfig(false)
	ethClient := evmclient.NewSimulatedBackendClient(t, b, testutils.SimulatedChainID)
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	lorm := logpoller.NewORM(testutils.SimulatedChainID, db, lggr, cfg)
	lp := logpoller.NewLogPoller(lorm, ethClient, lggr, 100*time.Millisecond, 1, 2, 2, 1000)
	require.NoError(t, lp.Start(ctx))
	t.Cleanup(func() { lp.Close() })

	t.Run("happy path", func(t *testing.T) {
		cp, err := NewConfigPoller(lggr, ethClient, lp, ocrAddress)
		require.NoError(t, err)
		// Should have no config to begin with.
		_, config, err := cp.LatestConfigDetails(testutils.Context(t))
		require.NoError(t, err)
		require.Equal(t, ocrtypes2.ConfigDigest{}, config)
		// Set the config
		contractConfig := setConfig(t, median.OffchainConfig{
			AlphaReportInfinite: false,
			AlphaReportPPB:      0,
			AlphaAcceptInfinite: true,
			AlphaAcceptPPB:      0,
			DeltaC:              10,
		}, ocrContract, user)
		b.Commit()
		latest, err := b.BlockByNumber(testutils.Context(t), nil)
		require.NoError(t, err)
		// Ensure we capture this config set log.
		require.NoError(t, lp.Replay(testutils.Context(t), latest.Number().Int64()-1))

		// Send blocks until we see the config updated.
		var configBlock uint64
		var digest [32]byte
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			b.Commit()
			configBlock, digest, err = cp.LatestConfigDetails(testutils.Context(t))
			require.NoError(t, err)
			return ocrtypes2.ConfigDigest{} != digest
		}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())

		// Assert the config returned is the one we configured.
		newConfig, err := cp.LatestConfig(testutils.Context(t), configBlock)
		require.NoError(t, err)
		// Note we don't check onchainConfig, as that is populated in the contract itself.
		assert.Equal(t, digest, [32]byte(newConfig.ConfigDigest))
		assert.Equal(t, contractConfig.Signers, newConfig.Signers)
		assert.Equal(t, contractConfig.Transmitters, newConfig.Transmitters)
		assert.Equal(t, contractConfig.F, newConfig.F)
		assert.Equal(t, contractConfig.OffchainConfigVersion, newConfig.OffchainConfigVersion)
		assert.Equal(t, contractConfig.OffchainConfig, newConfig.OffchainConfig)
	})

	ocrAddressPersistConfigEnabled, _, ocrContractPersistConfigEnabled, err := ocr2aggregator.DeployOCR2Aggregator(
		user,
		b,
		linkTokenAddress,
		big.NewInt(0),
		big.NewInt(10),
		accessAddress,
		accessAddress,
		9,
		"TEST",
		true,
	)
	require.NoError(t, err)
	b.Commit()

	t.Run("callIsConfigPersisted returns false", func(t *testing.T) {
		t.Run("when contract method missing, does not enable persistConfig", func(t *testing.T) {
			// Latest contract will always have the method, so use a random
			// unrelated contract (i.e. Link token) that doesn't implement the
			// method to test
			cp, err := newConfigPoller(lggr, ethClient, lp, linkTokenAddress)
			require.NoError(t, err)

			persistConfig, err := cp.callIsConfigPersisted(testutils.Context(t))
			require.NoError(t, err)
			assert.False(t, persistConfig)
		})
		t.Run("when contract method returns false, does not enable persistConfig", func(t *testing.T) {
			// Deployed test ocr contract has persistConfig=false
			cp, err := newConfigPoller(lggr, ethClient, lp, ocrAddress)
			require.NoError(t, err)

			persistConfig, err := cp.callIsConfigPersisted(testutils.Context(t))
			require.NoError(t, err)
			assert.False(t, persistConfig)
		})
	})

	t.Run("callIsConfigPersisted returns true", func(t *testing.T) {
		cp, err := newConfigPoller(lggr, ethClient, lp, ocrAddressPersistConfigEnabled)
		require.NoError(t, err)

		persistConfig, err := cp.callIsConfigPersisted(testutils.Context(t))
		require.NoError(t, err)
		assert.True(t, persistConfig)
	})

	t.Run("LatestConfigDetails, when logs have been pruned and persistConfig is true", func(t *testing.T) {
		// Give it a log poller that will never return logs
		mp := new(mocks.LogPoller)
		mp.On("RegisterFilter", mock.Anything).Return(nil)
		mp.On("LatestLogByEventSigWithConfs", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

		t.Run("if callLatestConfigDetails succeeds", func(t *testing.T) {
			cp, err := newConfigPoller(lggr, ethClient, mp, ocrAddressPersistConfigEnabled)
			require.NoError(t, err)
			cp.persistConfig.Store(true)

			t.Run("when config has not been set, returns zero values", func(t *testing.T) {
				changedInBlock, configDigest, err := cp.LatestConfigDetails(testutils.Context(t))
				require.NoError(t, err)

				assert.Equal(t, 0, int(changedInBlock))
				assert.Equal(t, ocrtypes.ConfigDigest{}, configDigest)
			})
			t.Run("when config has been set, returns config details", func(t *testing.T) {
				setConfig(t, median.OffchainConfig{
					AlphaReportInfinite: false,
					AlphaReportPPB:      0,
					AlphaAcceptInfinite: true,
					AlphaAcceptPPB:      0,
					DeltaC:              10,
				}, ocrContractPersistConfigEnabled, user)
				b.Commit()

				changedInBlock, configDigest, err := cp.LatestConfigDetails(testutils.Context(t))
				require.NoError(t, err)

				latest, err := b.BlockByNumber(testutils.Context(t), nil)
				require.NoError(t, err)

				onchainDetails, err := ocrContractPersistConfigEnabled.LatestConfigDetails(nil)
				require.NoError(t, err)

				assert.Equal(t, latest.Number().Int64(), int64(changedInBlock))
				assert.Equal(t, onchainDetails.ConfigDigest, [32]byte(configDigest))
			})
		})
		t.Run("returns error if callLatestConfigDetails fails", func(t *testing.T) {
			failingClient := new(evmClientMocks.Client)
			failingClient.On("ConfiguredChainID").Return(big.NewInt(42))
			failingClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("something exploded!"))
			cp, err := newConfigPoller(lggr, failingClient, mp, ocrAddressPersistConfigEnabled)
			require.NoError(t, err)
			cp.persistConfig.Store(true)

			_, _, err = cp.LatestConfigDetails(testutils.Context(t))
			assert.EqualError(t, err, "something exploded!")
		})
	})

	// deploy it again to reset to empty config
	ocrAddressPersistConfigEnabled, _, ocrContractPersistConfigEnabled, err = ocr2aggregator.DeployOCR2Aggregator(
		user,
		b,
		linkTokenAddress,
		big.NewInt(0),
		big.NewInt(10),
		accessAddress,
		accessAddress,
		9,
		"TEST",
		true,
	)
	require.NoError(t, err)
	b.Commit()

	t.Run("LatestConfig, when logs have been pruned and persistConfig is true", func(t *testing.T) {
		// Give it a log poller that will never return logs
		mp := new(mocks.LogPoller)
		mp.On("RegisterFilter", mock.Anything).Return(nil)
		mp.On("Logs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

		latest, err := b.BlockByNumber(testutils.Context(t), nil)
		require.NoError(t, err)
		blockNum := uint64(latest.Number().Int64())

		t.Run("if callLatestConfig succeeds", func(t *testing.T) {
			cp, err := newConfigPoller(lggr, ethClient, mp, ocrAddressPersistConfigEnabled)
			require.NoError(t, err)
			cp.persistConfig.Store(true)

			t.Run("when config has not been set, returns zero values", func(t *testing.T) {
				contractConfig, err := cp.LatestConfig(testutils.Context(t), blockNum)
				require.NoError(t, err)

				assert.Equal(t, ocrtypes.ConfigDigest{}, contractConfig.ConfigDigest)
			})
			t.Run("when config has been set, returns config details", func(t *testing.T) {
				contractConfig := setConfig(t, median.OffchainConfig{
					AlphaReportInfinite: false,
					AlphaReportPPB:      0,
					AlphaAcceptInfinite: true,
					AlphaAcceptPPB:      0,
					DeltaC:              10,
				}, ocrContractPersistConfigEnabled, user)
				b.Commit()
				blockNum++

				newConfig, err := cp.LatestConfig(testutils.Context(t), blockNum)
				require.NoError(t, err)

				onchainDetails, err := ocrContractPersistConfigEnabled.LatestConfigDetails(nil)
				require.NoError(t, err)

				assert.Equal(t, onchainDetails.ConfigDigest, [32]byte(newConfig.ConfigDigest))
				assert.Equal(t, contractConfig.Signers, newConfig.Signers)
				assert.Equal(t, contractConfig.Transmitters, newConfig.Transmitters)
				assert.Equal(t, contractConfig.F, newConfig.F)
				assert.Equal(t, contractConfig.OffchainConfigVersion, newConfig.OffchainConfigVersion)
				assert.Equal(t, contractConfig.OffchainConfig, newConfig.OffchainConfig)
			})
		})
		t.Run("returns error if callLatestConfig fails", func(t *testing.T) {
			failingClient := new(evmClientMocks.Client)
			failingClient.On("ConfiguredChainID").Return(big.NewInt(42))
			failingClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("something exploded!"))
			cp, err := newConfigPoller(lggr, failingClient, mp, ocrAddressPersistConfigEnabled)
			require.NoError(t, err)
			cp.persistConfig.Store(true)

			_, err = cp.LatestConfig(testutils.Context(t), blockNum)
			assert.EqualError(t, err, "something exploded!")
		})
	})
}

func setConfig(t *testing.T, pluginConfig median.OffchainConfig, ocrContract *ocr2aggregator.OCR2Aggregator, user *bind.TransactOpts) ocrtypes2.ContractConfig {
	// Create minimum number of nodes.
	var oracles []confighelper2.OracleIdentityExtra
	for i := 0; i < 4; i++ {
		oracles = append(oracles, confighelper2.OracleIdentityExtra{
			OracleIdentity: confighelper2.OracleIdentity{
				OnchainPublicKey:  utils.RandomAddress().Bytes(),
				TransmitAccount:   ocrtypes2.Account(fmt.Sprintf("0x%x", utils.RandomAddress().Bytes())),
				OffchainPublicKey: utils.RandomBytes32(),
				PeerID:            utils.MustNewPeerID(),
			},
			ConfigEncryptionPublicKey: utils.RandomBytes32(),
		})
	}
	// Change the offramp config
	signers, transmitters, threshold, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper2.ContractSetConfigArgsForTests(
		2*time.Second,        // deltaProgress
		1*time.Second,        // deltaResend
		1*time.Second,        // deltaRound
		500*time.Millisecond, // deltaGrace
		2*time.Second,        // deltaStage
		3,
		[]int{1, 1, 1, 1},
		oracles,
		pluginConfig.Encode(),
		50*time.Millisecond,
		50*time.Millisecond,
		50*time.Millisecond,
		50*time.Millisecond,
		50*time.Millisecond,
		1, // faults
		nil,
	)
	require.NoError(t, err)
	signerAddresses, err := OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)
	transmitterAddresses, err := AccountToAddress(transmitters)
	require.NoError(t, err)
	_, err = ocrContract.SetConfig(user, signerAddresses, transmitterAddresses, threshold, onchainConfig, offchainConfigVersion, offchainConfig)
	require.NoError(t, err)
	return ocrtypes2.ContractConfig{
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     threshold,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}
}

func ptr[T any](v T) *T { return &v }
