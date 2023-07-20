// Package testsetups compresses common test setups and more complicated setups like performance and chaos tests.
package testsetups

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/ethereum"
	"github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver"
	mockservercfg "github.com/smartcontractkit/chainlink-env/pkg/helm/mockserver-cfg"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	reportModel "github.com/smartcontractkit/chainlink-testing-framework/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/config"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/networks"
	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
)

// OCRSoakTest defines a typical OCR soak test
type OCRSoakTest struct {
	Inputs                *OCRSoakTestInputs
	TestReporter          testreporters.OCRSoakTestReporter
	OperatorForwarderFlow bool

	t               *testing.T
	testEnvironment *environment.Environment
	log             zerolog.Logger
	bootstrapNode   *client.Chainlink
	workerNodes     []*client.Chainlink
	chainClient     blockchain.EVMClient
	mockServer      *ctfClient.MockserverClient
	filterQuery     geth.FilterQuery

	ocrRoundStates []*testreporters.OCRRoundState
	rpcIssues      []*testreporters.RPCIssue

	ocrInstances   []contracts.OffchainAggregator
	ocrInstanceMap map[string]contracts.OffchainAggregator // address : instance
}

// OCRSoakTestInputs define required inputs to run an OCR soak test
type OCRSoakTestInputs struct {
	TestDuration            time.Duration `envconfig:"TEST_DURATION" default:"15m"`         // How long to run the test for
	NumberOfContracts       int           `envconfig:"NUMBER_CONTRACTS" default:"2"`        // Number of OCR contracts to launch
	ChainlinkNodeFunding    float64       `envconfig:"CHAINLINK_NODE_FUNDING" default:".1"` // Amount of native currency to fund each chainlink node with
	bigChainlinkNodeFunding *big.Float    // Convenience conversions for funding
	TimeBetweenRounds       time.Duration `envconfig:"TIME_BETWEEN_ROUNDS" default:"1m"` // How long to wait before starting a new round; controls frequency of rounds
}

func (i OCRSoakTestInputs) setForRemoteRunner() {
	os.Setenv("TEST_OCR_TEST_DURATION", i.TestDuration.String())
	os.Setenv("TEST_OCR_NUMBER_CONTRACTS", fmt.Sprint(i.NumberOfContracts))
	os.Setenv("TEST_OCR_CHAINLINK_NODE_FUNDING", strconv.FormatFloat(i.ChainlinkNodeFunding, 'f', -1, 64))
	os.Setenv("TEST_OCR_TIME_BETWEEN_ROUNDS", i.TimeBetweenRounds.String())

	selectedNetworks := strings.Split(os.Getenv("SELECTED_NETWORKS"), ",")
	for _, networkPrefix := range selectedNetworks {
		urlEnv := fmt.Sprintf("%s_URLS", networkPrefix)
		httpEnv := fmt.Sprintf("%s_HTTP_URLS", networkPrefix)
		os.Setenv(fmt.Sprintf("TEST_%s", urlEnv), os.Getenv(urlEnv))
		os.Setenv(fmt.Sprintf("TEST_%s", httpEnv), os.Getenv(httpEnv))
	}
}

// NewOCRSoakTest creates a new OCR soak test to setup and run
func NewOCRSoakTest(t *testing.T, forwarderFlow bool) (*OCRSoakTest, error) {
	var testInputs OCRSoakTestInputs
	err := envconfig.Process("OCR", &testInputs)
	if err != nil {
		return nil, err
	}
	testInputs.setForRemoteRunner()

	test := &OCRSoakTest{
		Inputs:                &testInputs,
		OperatorForwarderFlow: forwarderFlow,
		TestReporter: testreporters.OCRSoakTestReporter{
			TestDuration: testInputs.TestDuration,
		},
		t:              t,
		log:            utils.GetTestLogger(t),
		ocrRoundStates: make([]*testreporters.OCRRoundState, 0),
		ocrInstanceMap: make(map[string]contracts.OffchainAggregator),
	}
	test.ensureInputValues()
	return test, nil
}

// DeployEnvironment deploys the test environment, starting all Chainlink nodes and other components for the test
func (o *OCRSoakTest) DeployEnvironment(customChainlinkNetworkTOML string) {
	network := networks.SelectedNetwork // Environment currently being used to soak test on
	nsPre := "soak-ocr-"
	if o.OperatorForwarderFlow {
		nsPre = fmt.Sprintf("%sforwarder-", nsPre)
	}
	nsPre = fmt.Sprintf("%s%s", nsPre, strings.ReplaceAll(strings.ToLower(network.Name), " ", "-"))
	baseEnvironmentConfig := &environment.Config{
		TTL:             time.Hour * 720, // 30 days,
		NamespacePrefix: nsPre,
		Test:            o.t,
	}

	cd, err := chainlink.NewDeployment(6, map[string]any{
		"toml": client.AddNetworkDetailedConfig(config.BaseOCRP2PV1Config, customChainlinkNetworkTOML, network),
		"db": map[string]any{
			"stateful": true, // stateful DB by default for soak tests
		},
	})
	require.NoError(o.t, err, "Error creating chainlink deployment")
	testEnvironment := environment.New(baseEnvironmentConfig).
		AddHelm(mockservercfg.New(nil)).
		AddHelm(mockserver.New(nil)).
		AddHelm(ethereum.New(&ethereum.Props{
			NetworkName: network.Name,
			Simulated:   network.Simulated,
			WsURLs:      network.URLs,
		})).
		AddHelmCharts(cd)
	err = testEnvironment.Run()
	require.NoError(o.t, err, "Error launching test environment")
	o.testEnvironment = testEnvironment
}

// LoadEnvironment loads an existing test environment using the provided URLs
func (o *OCRSoakTest) LoadEnvironment(chainlinkURLs []string, chainURL, mockServerURL string) {
	var (
		network = networks.SelectedNetwork
		err     error
	)
	o.chainClient, err = blockchain.ConnectEVMClient(network)
	require.NoError(o.t, err, "Error connecting to EVM client")
	chainlinkNodes, err := client.ConnectChainlinkNodeURLs(chainlinkURLs)
	require.NoError(o.t, err, "Error connecting to chainlink nodes")
	o.bootstrapNode, o.workerNodes = chainlinkNodes[0], chainlinkNodes[1:]
	o.mockServer, err = ctfClient.ConnectMockServerURL(mockServerURL)
	require.NoError(o.t, err, "Error connecting to mockserver")
}

// Environment returns the full K8s test environment
func (o *OCRSoakTest) Environment() *environment.Environment {
	return o.testEnvironment
}

func (o *OCRSoakTest) Setup() {
	var (
		err     error
		network = networks.SelectedNetwork
	)
	// Environment currently being used to soak test on
	// Make connections to soak test resources
	o.chainClient, err = blockchain.NewEVMClient(network, o.testEnvironment)
	require.NoError(o.t, err, "Error creating EVM client")
	contractDeployer, err := contracts.NewContractDeployer(o.chainClient)
	require.NoError(o.t, err, "Unable to create contract deployer")
	require.NotNil(o.t, contractDeployer, "Contract deployer shouldn't be nil")
	nodes, err := client.ConnectChainlinkNodes(o.testEnvironment)
	require.NoError(o.t, err, "Connecting to chainlink nodes shouldn't fail")
	o.bootstrapNode, o.workerNodes = nodes[0], nodes[1:]
	o.mockServer, err = ctfClient.ConnectMockServer(o.testEnvironment)
	require.NoError(o.t, err, "Creating mockserver clients shouldn't fail")
	o.chainClient.ParallelTransactions(true)
	// Deploy LINK
	linkTokenContract, err := contractDeployer.DeployLinkTokenContract()
	require.NoError(o.t, err, "Deploying Link Token Contract shouldn't fail")

	// Fund Chainlink nodes, excluding the bootstrap node
	err = actions.FundChainlinkNodes(o.workerNodes, o.chainClient, o.Inputs.bigChainlinkNodeFunding)
	require.NoError(o.t, err, "Error funding Chainlink nodes")

	if o.OperatorForwarderFlow {
		contractLoader, err := contracts.NewContractLoader(o.chainClient)
		require.NoError(o.t, err, "Loading contracts shouldn't fail")

		operators, authorizedForwarders, _ := actions.DeployForwarderContracts(
			o.t, contractDeployer, linkTokenContract, o.chainClient, len(o.workerNodes),
		)
		forwarderNodesAddresses, err := actions.ChainlinkNodeAddresses(o.workerNodes)
		require.NoError(o.t, err, "Retreiving on-chain wallet addresses for chainlink nodes shouldn't fail")
		for i := range o.workerNodes {
			actions.AcceptAuthorizedReceiversOperator(
				o.t, operators[i], authorizedForwarders[i], []common.Address{forwarderNodesAddresses[i]}, o.chainClient, contractLoader,
			)
			require.NoError(o.t, err, "Accepting Authorize Receivers on Operator shouldn't fail")
			actions.TrackForwarder(o.t, o.chainClient, authorizedForwarders[i], o.workerNodes[i])
			err = o.chainClient.WaitForEvents()
		}

		o.ocrInstances = actions.DeployOCRContractsForwarderFlow(
			o.t,
			o.Inputs.NumberOfContracts,
			linkTokenContract,
			contractDeployer,
			o.workerNodes,
			authorizedForwarders,
			o.chainClient,
		)
	} else {
		o.ocrInstances, err = actions.DeployOCRContracts(
			o.Inputs.NumberOfContracts,
			linkTokenContract,
			contractDeployer,
			o.bootstrapNode,
			o.workerNodes,
			o.chainClient,
		)
		require.NoError(o.t, err)
	}

	err = o.chainClient.WaitForEvents()
	require.NoError(o.t, err, "Error waiting for OCR contracts to be deployed")
	for _, ocrInstance := range o.ocrInstances {
		o.ocrInstanceMap[ocrInstance.Address()] = ocrInstance
	}
	o.log.Info().Msg("OCR Soak Test Setup Complete")
}

// Run starts the OCR soak test
func (o *OCRSoakTest) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	latestBlockNum, err := o.chainClient.LatestBlockNumber(ctx)
	cancel()
	require.NoError(o.t, err, "Error getting current block number")

	ocrAddresses := make([]common.Address, len(o.ocrInstances))
	for i, ocrInstance := range o.ocrInstances {
		ocrAddresses[i] = common.HexToAddress(ocrInstance.Address())
	}
	contractABI, err := offchainaggregator.OffchainAggregatorMetaData.GetAbi()
	require.NoError(o.t, err, "Error retrieving OCR contract ABI")
	o.filterQuery = geth.FilterQuery{
		Addresses: ocrAddresses,
		Topics:    [][]common.Hash{{contractABI.Events["AnswerUpdated"].ID}},
		FromBlock: big.NewInt(0).SetUint64(latestBlockNum),
	}

	if o.OperatorForwarderFlow {
		actions.CreateOCRJobsWithForwarder(o.t, o.ocrInstances, o.bootstrapNode, o.workerNodes, 5, o.mockServer)
	} else {
		err := actions.CreateOCRJobs(o.ocrInstances, o.bootstrapNode, o.workerNodes, 5, o.mockServer)
		require.NoError(o.t, err, "Error creating OCR jobs")
	}

	o.log.Info().
		Str("Test Duration", o.Inputs.TestDuration.Truncate(time.Second).String()).
		Int("Number of OCR Contracts", len(o.ocrInstances)).
		Msg("Starting OCR Soak Test")

	testDuration := time.After(o.Inputs.TestDuration)

	// *********************
	// ***** Test Loop *****
	// *********************
	interruption := make(chan os.Signal, 1)
	signal.Notify(interruption, os.Kill, os.Interrupt, syscall.SIGTERM)
	lastValue := 0
	newRoundTrigger := time.NewTimer(0) // Want to trigger a new round ASAP
	defer newRoundTrigger.Stop()
	err = o.observeOCREvents()
	require.NoError(o.t, err, "Error subscribing to OCR events")

testLoop:
	for {
		select {
		case <-interruption:
			o.log.Info().Msg("Test interrupted, shutting down")
			os.Exit(1)
		case <-testDuration:
			break testLoop
		case <-newRoundTrigger.C:
			newValue := rand.Intn(256) + 1 // #nosec G404 - not everything needs to be cryptographically secure
			for newValue == lastValue {
				newValue = rand.Intn(256) + 1 // #nosec G404 - kudos to you if you actually find a way to exploit this
			}
			lastValue = newValue
			err := o.triggerNewRound(newValue)

			timerReset := o.Inputs.TimeBetweenRounds
			if err != nil {
				timerReset = time.Second * 5
				o.log.Error().Err(err).
					Int("Seconds Waiting", 5).
					Msg("Error triggering new round, waiting and trying again. Possible connection issues with mockserver")
			}
			newRoundTrigger.Reset(timerReset)
		case t := <-o.chainClient.ConnectionIssue():
			o.rpcIssues = append(o.rpcIssues, &testreporters.RPCIssue{
				StartTime: t,
				Message:   "RPC Connection Lost",
			})
		case t := <-o.chainClient.ConnectionRestored():
			o.rpcIssues = append(o.rpcIssues, &testreporters.RPCIssue{
				StartTime: t,
				Message:   "RPC Connection Restored",
			})
		}
	}

	o.log.Info().Msg("Test Complete, collecting on-chain events to be collected")
	// Keep trying to collect events until we get them, no exceptions
	timeout := time.Second * 5
	err = o.collectEvents(timeout)
	for err != nil {
		timeout *= 2
		err = o.collectEvents(timeout)
	}
	o.TestReporter.RecordEvents(o.ocrRoundStates, o.rpcIssues)
}

func (o *OCRSoakTest) SaveState() {

}

func (o *OCRSoakTest) LoadState() {

}

func (o *OCRSoakTest) Resume() {

}

// Networks returns the networks that the test is running on
func (o *OCRSoakTest) TearDownVals(t *testing.T) (
	*testing.T,
	*environment.Environment,
	[]*client.Chainlink,
	reportModel.TestReporter,
	blockchain.EVMClient,
) {
	return t, o.testEnvironment, append(o.workerNodes, o.bootstrapNode), &o.TestReporter, o.chainClient
}

// *********************
// Recovery if the test is shut-down/rebalanced by K8s
// *********************

// OCRSoakTestState contains all the info needed by the test to recover from a K8s rebalance, assuming the test was in a running state
type OCRSoakTestState struct {
	OCRRoundStates       []*testreporters.OCRRoundState `toml:"ocrRoundStates"`
	RPCIssues            []*testreporters.RPCIssue      `toml:"rpcIssues"`
	TimeRunning          time.Duration                  `toml:"timeRunning"`
	TestDuration         time.Duration                  `toml:"testDuration"`
	OCRContractAddresses []string                       `toml:"ocrContractAddresses"`

	BootStrapNodeURL string   `toml:"bootstrapNodeURL"`
	WorkerNodeURLs   []string `toml:"workerNodeURLs"`
	ChainURL         string   `toml:"chainURL"`
	MockServerURL    string   `toml:"mockServerURL"`
}

// *********************
// ****** Helpers ******
// *********************

// observeOCREvents subscribes to OCR events and logs them to the test logger
// WARNING: Should only be used for observation and logging. This is not a reliable way to collect events.
func (o *OCRSoakTest) observeOCREvents() error {
	eventLogs := make(chan types.Log)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	eventSub, err := o.chainClient.SubscribeFilterLogs(ctx, o.filterQuery, eventLogs)
	if err != nil {
		return err
	}

	go func() {
		defer cancel()
		for {
			select {
			case event := <-eventLogs:
				answerUpdated, err := o.ocrInstances[0].ParseEventAnswerUpdated(event)
				if err != nil {
					log.Warn().
						Err(err).
						Str("Address", event.Address.Hex()).
						Uint64("Block Number", event.BlockNumber).
						Msg("Error parsing event as AnswerUpdated")
					continue
				}
				o.log.Info().
					Str("Address", event.Address.Hex()).
					Uint64("Block Number", event.BlockNumber).
					Uint64("Round ID", answerUpdated.RoundId.Uint64()).
					Int64("Answer", answerUpdated.Current.Int64()).
					Msg("Answer Updated Event")
			case err = <-eventSub.Err():
				for err != nil {
					o.log.Trace().
						Err(err).
						Interface("Query", o.filterQuery).
						Msg("Error while subscribed to OCR Logs. Resubscribing")
					ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
					eventSub, err = o.chainClient.SubscribeFilterLogs(ctx, o.filterQuery, eventLogs)
				}
			}
		}
	}()

	return nil
}

// triggers a new OCR round by setting a new mock adapter value
func (o *OCRSoakTest) triggerNewRound(newValue int) error {
	if len(o.ocrRoundStates) > 0 {
		o.ocrRoundStates[len(o.ocrRoundStates)-1].EndTime = time.Now()
	}

	err := actions.SetAllAdapterResponsesToTheSameValue(newValue, o.ocrInstances, o.workerNodes, o.mockServer)
	if err != nil {
		return err
	}

	expectedState := &testreporters.OCRRoundState{
		StartTime:   time.Now(),
		Answer:      int64(newValue),
		FoundEvents: make(map[string][]*testreporters.FoundEvent),
	}
	for _, ocrInstance := range o.ocrInstances {
		expectedState.FoundEvents[ocrInstance.Address()] = make([]*testreporters.FoundEvent, 0)
	}
	o.ocrRoundStates = append(o.ocrRoundStates, expectedState)
	o.log.Info().
		Int("Value", newValue).
		Msg("Starting a New OCR Round")
	return nil
}

func (o *OCRSoakTest) collectEvents(timeout time.Duration) error {
	start := time.Now()
	o.ocrRoundStates[len(o.ocrRoundStates)-1].EndTime = start // Set end time for last expected event
	o.log.Info().Msg("Collecting on-chain events")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	contractEvents, err := o.chainClient.FilterLogs(ctx, o.filterQuery)
	if err != nil {
		log.Error().
			Err(err).
			Str("Time", time.Since(start).String()).
			Msg("Error collecting on-chain events")
		return err
	}

	sortedFoundEvents := make([]*testreporters.FoundEvent, 0)
	for _, event := range contractEvents {
		answerUpdated, err := o.ocrInstances[0].ParseEventAnswerUpdated(event)
		if err != nil {
			log.Error().
				Err(err).
				Str("Time", time.Since(start).String()).
				Msg("Error collecting on-chain events")
			return err
		}
		sortedFoundEvents = append(sortedFoundEvents, &testreporters.FoundEvent{
			StartTime:   time.Unix(answerUpdated.UpdatedAt.Int64(), 0),
			Address:     event.Address.Hex(),
			Answer:      answerUpdated.Current.Int64(),
			RoundID:     answerUpdated.RoundId.Uint64(),
			BlockNumber: event.BlockNumber,
		})
	}

	// Sort our events by time to make sure they are in order (don't trust RPCs)
	sort.Slice(sortedFoundEvents, func(i, j int) bool {
		return sortedFoundEvents[i].StartTime.Before(sortedFoundEvents[j].StartTime)
	})

	// Now match each found event with the expected event time frame
	expectedIndex := 0
	for _, event := range sortedFoundEvents {
		if !event.StartTime.Before(o.ocrRoundStates[expectedIndex].EndTime) {
			expectedIndex++
			if expectedIndex >= len(o.ocrRoundStates) {
				o.log.Warn().
					Str("Event Time", event.StartTime.String()).
					Str("Expected End Time", o.ocrRoundStates[expectedIndex].EndTime.String()).
					Msg("Found events after last expected end time, adding event to that final report, things might be weird")
			}
		}
		o.ocrRoundStates[expectedIndex].FoundEvents[event.Address] = append(o.ocrRoundStates[expectedIndex].FoundEvents[event.Address], event)
		o.ocrRoundStates[expectedIndex].TimeLineEvents = append(o.ocrRoundStates[expectedIndex].TimeLineEvents, event)
	}

	o.log.Info().
		Str("Time", time.Since(start).String()).
		Msg("Collected on-chain events")
	return nil
}

// ensureValues ensures that all values needed to run the test are present
func (o *OCRSoakTest) ensureInputValues() error {
	inputs := o.Inputs
	if inputs.NumberOfContracts <= 0 {
		return fmt.Errorf("Number of OCR contracts must be greater than 0, found %d", inputs.NumberOfContracts)
	}
	if inputs.ChainlinkNodeFunding <= 0 {
		return fmt.Errorf("Chainlink node funding must be greater than 0, found %f", inputs.ChainlinkNodeFunding)
	}
	if inputs.TestDuration <= time.Minute {
		return fmt.Errorf("Test duration must be greater than 1 minute, found %s", inputs.TestDuration.String())
	}
	if inputs.TimeBetweenRounds >= time.Hour {
		return fmt.Errorf("Time between rounds must be less than 1 hour, found %s", inputs.TimeBetweenRounds.String())
	}
	o.Inputs.bigChainlinkNodeFunding = big.NewFloat(inputs.ChainlinkNodeFunding)
	return nil
}
