package actions

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

// This actions file often returns functions, rather than just values. These are used as common test helpers, and are
// handy to have returning as functions so that Ginkgo can use them in an aesthetically pleasing way.

// DeployOCRContracts deploys and funds a certain number of offchain aggregator contracts
func DeployOCRContracts(
	numberOfContracts int,
	linkTokenContract contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	bootstrapNode *client.Chainlink,
	workerNodes []*client.Chainlink,
	client blockchain.EVMClient,
) ([]contracts.OffchainAggregator, error) {
	// Deploy contracts
	var ocrInstances []contracts.OffchainAggregator
	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		ocrInstance, err := contractDeployer.DeployOffChainAggregator(
			linkTokenContract.Address(),
			contracts.DefaultOffChainAggregatorOptions(),
		)
		if err != nil {
			return nil, fmt.Errorf("OCR instance deployment have failed: %w", err)
		}
		ocrInstances = append(ocrInstances, ocrInstance)
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			if err != nil {
				return nil, fmt.Errorf("failed to wait for OCR contract deployments: %w", err)
			}
		}
	}
	err := client.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("error waiting for OCR contract deployments: %w", err)
	}

	// Gather transmitter and address payees
	var transmitters, payees []string
	for _, node := range workerNodes {
		addr, err := node.PrimaryEthAddress()
		if err != nil {
			return nil, fmt.Errorf("error getting node's primary ETH address: %w", err)
		}
		transmitters = append(transmitters, addr)
		payees = append(payees, client.GetDefaultWallet().Address())
	}

	// Set Payees
	for contractCount, ocrInstance := range ocrInstances {
		err = ocrInstance.SetPayees(transmitters, payees)
		if err != nil {
			return nil, fmt.Errorf("error settings OCR payees: %w", err)
		}
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			if err != nil {
				return nil, fmt.Errorf("failed to wait for setting OCR payees: %w", err)
			}
		}
	}
	err = client.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("error waiting for OCR contracts to set payees and transmitters: %w", err)
	}

	// Set Config
	transmitterAddresses, err := ChainlinkNodeAddresses(workerNodes)
	if err != nil {
		return nil, fmt.Errorf("getting node common addresses should not fail: %w", err)
	}
	for contractCount, ocrInstance := range ocrInstances {
		// Exclude the first node, which will be used as a bootstrapper
		err = ocrInstance.SetConfig(
			workerNodes,
			contracts.DefaultOffChainAggregatorConfig(len(workerNodes)),
			transmitterAddresses,
		)
		if err != nil {
			return nil, fmt.Errorf("error setting OCR config for contract '%s': %w", ocrInstance.Address(), err)
		}
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			if err != nil {
				return nil, fmt.Errorf("failed to wait for setting OCR config: %w", err)
			}
		}
	}
	err = client.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("error waiting for OCR contracts to set config: %w", err)
	}
	return ocrInstances, nil
}

// DeployOCRContractsForwarderFlow deploys and funds a certain number of offchain
// aggregator contracts with forwarders as effectiveTransmitters
func DeployOCRContractsForwarderFlow(
	t *testing.T,
	numberOfContracts int,
	linkTokenContract contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	workerNodes []*client.Chainlink,
	forwarderAddresses []common.Address,
	client blockchain.EVMClient,
) []contracts.OffchainAggregator {
	// Deploy contracts
	var ocrInstances []contracts.OffchainAggregator
	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		ocrInstance, err := contractDeployer.DeployOffChainAggregator(
			linkTokenContract.Address(),
			contracts.DefaultOffChainAggregatorOptions(),
		)
		require.NoError(t, err, "Deploying OCR instance %d shouldn't fail", contractCount+1)
		ocrInstances = append(ocrInstances, ocrInstance)
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			require.NoError(t, err, "Failed to wait for OCR Contract deployments")
		}
	}
	err := client.WaitForEvents()
	require.NoError(t, err, "Error waiting for OCR contract deployments")

	// Gather transmitter and address payees
	var transmitters, payees []string
	for _, forwarderCommonAddress := range forwarderAddresses {
		forwarderAddress := forwarderCommonAddress.Hex()
		transmitters = append(transmitters, forwarderAddress)
		payees = append(payees, client.GetDefaultWallet().Address())
	}

	// Set Payees
	for contractCount, ocrInstance := range ocrInstances {
		err = ocrInstance.SetPayees(transmitters, payees)
		require.NoError(t, err, "Error setting OCR payees")
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			require.NoError(t, err, "Failed to wait for setting OCR payees")
		}
	}
	err = client.WaitForEvents()
	require.NoError(t, err, "Error waiting for OCR contracts to set payees and transmitters")

	// Set Config
	for contractCount, ocrInstance := range ocrInstances {
		// Exclude the first node, which will be used as a bootstrapper
		err = ocrInstance.SetConfig(
			workerNodes,
			contracts.DefaultOffChainAggregatorConfig(len(workerNodes)),
			forwarderAddresses,
		)
		require.NoError(t, err, "Error setting OCR config for contract '%d'", ocrInstance.Address())
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			require.NoError(t, err, "Failed to wait for setting OCR config")
		}
	}
	err = client.WaitForEvents()
	require.NoError(t, err, "Error waiting for OCR contracts to set config")
	return ocrInstances
}

// CreateOCRJobs bootstraps the first node and to the other nodes sends ocr jobs that
// read from different adapters, to be used in combination with SetAdapterResponses
func CreateOCRJobs(
	ocrInstances []contracts.OffchainAggregator,
	bootstrapNode *client.Chainlink,
	workerNodes []*client.Chainlink,
	mockValue int,
	mockserver *ctfClient.MockserverClient,
) error {
	for _, ocrInstance := range ocrInstances {
		bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
		if err != nil {
			return fmt.Errorf("reading P2P keys from bootstrap node have failed: %w", err)
		}
		bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
		bootstrapSpec := &client.OCRBootstrapJobSpec{
			Name:            fmt.Sprintf("bootstrap-%s", uuid.New().String()),
			ContractAddress: ocrInstance.Address(),
			P2PPeerID:       bootstrapP2PId,
			IsBootstrapPeer: true,
		}
		_, err = bootstrapNode.MustCreateJob(bootstrapSpec)
		if err != nil {
			return fmt.Errorf("creating bootstrap job have failed: %w", err)
		}

		for _, node := range workerNodes {
			nodeP2PIds, err := node.MustReadP2PKeys()
			if err != nil {
				return fmt.Errorf("reading P2P keys from OCR node have failed: %w", err)
			}
			nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
			nodeTransmitterAddress, err := node.PrimaryEthAddress()
			if err != nil {
				return fmt.Errorf("getting primary ETH address from OCR node have failed: %w", err)
			}
			nodeOCRKeys, err := node.MustReadOCRKeys()
			if err != nil {
				return fmt.Errorf("getting OCR keys from OCR node have failed: %w", err)
			}
			nodeOCRKeyId := nodeOCRKeys.Data[0].ID

			nodeContractPairID, err := BuildNodeContractPairID(node, ocrInstance)
			if err != nil {
				return err
			}
			bta := &client.BridgeTypeAttributes{
				Name: nodeContractPairID,
				URL:  fmt.Sprintf("%s/%s", mockserver.Config.ClusterURL, strings.TrimPrefix(nodeContractPairID, "/")),
			}
			err = SetAdapterResponse(mockValue, ocrInstance, node, mockserver)
			if err != nil {
				return fmt.Errorf("setting adapter response for OCR node failed: %w", err)
			}
			err = node.MustCreateBridge(bta)
			if err != nil {
				return fmt.Errorf("creating bridge job have failed: %w", err)
			}

			ocrSpec := &client.OCRTaskJobSpec{
				ContractAddress:    ocrInstance.Address(),
				P2PPeerID:          nodeP2PId,
				P2PBootstrapPeers:  []*client.Chainlink{bootstrapNode},
				KeyBundleID:        nodeOCRKeyId,
				TransmitterAddress: nodeTransmitterAddress,
				ObservationSource:  client.ObservationSourceSpecBridge(bta),
			}
			_, err = node.MustCreateJob(ocrSpec)
			if err != nil {
				return fmt.Errorf("creating OCR task job on OCR node have failed: %w", err)
			}
		}
	}
	return nil
}

// CreateOCRJobsWithForwarder bootstraps the first node and to the other nodes sends ocr jobs that
// read from different adapters, to be used in combination with SetAdapterResponses
func CreateOCRJobsWithForwarder(
	t *testing.T,
	ocrInstances []contracts.OffchainAggregator,
	bootstrapNode *client.Chainlink,
	workerNodes []*client.Chainlink,
	mockValue int,
	mockserver *ctfClient.MockserverClient,
) {
	for _, ocrInstance := range ocrInstances {
		bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
		require.NoError(t, err, "Shouldn't fail reading P2P keys from bootstrap node")
		bootstrapP2PId := bootstrapP2PIds.Data[0].Attributes.PeerID
		bootstrapSpec := &client.OCRBootstrapJobSpec{
			Name:            fmt.Sprintf("bootstrap-%s", uuid.New().String()),
			ContractAddress: ocrInstance.Address(),
			P2PPeerID:       bootstrapP2PId,
			IsBootstrapPeer: true,
		}
		_, err = bootstrapNode.MustCreateJob(bootstrapSpec)
		require.NoError(t, err, "Shouldn't fail creating bootstrap job on bootstrap node")

		for nodeIndex, node := range workerNodes {
			nodeP2PIds, err := node.MustReadP2PKeys()
			require.NoError(t, err, "Shouldn't fail reading P2P keys from OCR node %d", nodeIndex+1)
			nodeP2PId := nodeP2PIds.Data[0].Attributes.PeerID
			nodeTransmitterAddress, err := node.PrimaryEthAddress()
			require.NoError(t, err, "Shouldn't fail getting primary ETH address from OCR node %d", nodeIndex+1)
			nodeOCRKeys, err := node.MustReadOCRKeys()
			require.NoError(t, err, "Shouldn't fail getting OCR keys from OCR node %d", nodeIndex+1)
			nodeOCRKeyId := nodeOCRKeys.Data[0].ID

			nodeContractPairID, err := BuildNodeContractPairID(node, ocrInstance)
			require.NoError(t, err, "Failed building node contract pair ID for mockserver")
			bta := &client.BridgeTypeAttributes{
				Name: nodeContractPairID,
				URL:  fmt.Sprintf("%s/%s", mockserver.Config.ClusterURL, strings.TrimPrefix(nodeContractPairID, "/")),
			}
			err = SetAdapterResponse(mockValue, ocrInstance, node, mockserver)
			require.NoError(t, err, "Failed setting adapter responses for node %d", nodeIndex+1)
			err = node.MustCreateBridge(bta)
			require.NoError(t, err, "Failed creating bridge on OCR node %d", nodeIndex+1)

			ocrSpec := &client.OCRTaskJobSpec{
				ContractAddress:    ocrInstance.Address(),
				P2PPeerID:          nodeP2PId,
				P2PBootstrapPeers:  []*client.Chainlink{bootstrapNode},
				KeyBundleID:        nodeOCRKeyId,
				TransmitterAddress: nodeTransmitterAddress,
				ObservationSource:  client.ObservationSourceSpecBridge(bta),
				ForwardingAllowed:  true,
			}
			_, err = node.MustCreateJob(ocrSpec)
			require.NoError(t, err, "Shouldn't fail creating OCR Task job on OCR node %d", nodeIndex+1)
		}
	}
}

// StartNewRound requests a new round from the ocr contracts and waits for confirmation
func StartNewRound(
	roundNumber int64,
	ocrInstances []contracts.OffchainAggregator,
	client blockchain.EVMClient,
) error {
	for i := 0; i < len(ocrInstances); i++ {
		err := ocrInstances[i].RequestNewRound()
		if err != nil {
			return fmt.Errorf("requesting new OCR round %d have failed: %w", i+1, err)
		}
		ocrRound := contracts.NewOffchainAggregatorRoundConfirmer(ocrInstances[i], big.NewInt(roundNumber), client.GetNetworkConfig().Timeout.Duration)
		client.AddHeaderEventSubscription(ocrInstances[i].Address(), ocrRound)
		err = client.WaitForEvents()
		if err != nil {
			return fmt.Errorf("failed to wait for event subscriptions of OCR instance %d: %w", i+1, err)
		}
	}
	return nil
}

// SetAdapterResponse sets a single adapter response that correlates with an ocr contract and a chainlink node
func SetAdapterResponse(
	response int,
	ocrInstance contracts.OffchainAggregator,
	chainlinkNode *client.Chainlink,
	mockserver *ctfClient.MockserverClient,
) error {
	nodeContractPairID, err := BuildNodeContractPairID(chainlinkNode, ocrInstance)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/%s", nodeContractPairID)
	err = mockserver.SetValuePath(path, response)
	if err != nil {
		return fmt.Errorf("setting mockserver value path failed: %w", err)
	}
	return nil
}

// SetAllAdapterResponsesToTheSameValue sets the mock responses in mockserver that are read by chainlink nodes
// to simulate different adapters. This sets all adapter responses for each node and contract to the same response
func SetAllAdapterResponsesToTheSameValue(
	response int,
	ocrInstances []contracts.OffchainAggregator,
	chainlinkNodes []*client.Chainlink,
	mockserver *ctfClient.MockserverClient,
) error {
	eg := &errgroup.Group{}
	for _, o := range ocrInstances {
		ocrInstance := o
		for _, n := range chainlinkNodes {
			node := n
			eg.Go(func() error {
				return SetAdapterResponse(response, ocrInstance, node, mockserver)
			})
		}
	}
	return eg.Wait()
}

// SetAllAdapterResponsesToDifferentValues sets the mock responses in mockserver that are read by chainlink nodes
// to simulate different adapters. This sets all adapter responses for each node and contract to different responses
func SetAllAdapterResponsesToDifferentValues(
	t *testing.T,
	responses []int,
	ocrInstances []contracts.OffchainAggregator,
	chainlinkNodes []*client.Chainlink,
	mockserver *ctfClient.MockserverClient,
) {
	require.Equal(t, len(chainlinkNodes)-1, len(responses),
		"Amount of answers %d should be equal to the amount of Chainlink nodes - 1 for the bootstrap %d", len(responses), len(chainlinkNodes)-1)
	for _, ocrInstance := range ocrInstances {
		for nodeIndex := 1; nodeIndex < len(chainlinkNodes); nodeIndex++ {
			err := SetAdapterResponse(responses[nodeIndex-1], ocrInstance, chainlinkNodes[nodeIndex], mockserver)
			require.NoError(t, err)
		}
	}
}

// BuildNodeContractPairID builds a UUID based on a related pair of a Chainlink node and OCR contract
func BuildNodeContractPairID(node *client.Chainlink, ocrInstance contracts.OffchainAggregator) (string, error) {
	if node == nil {
		return "", fmt.Errorf("chainlink node is nil")
	}
	if ocrInstance == nil {
		return "", fmt.Errorf("OCR Instance is nil")
	}
	nodeAddress, err := node.PrimaryEthAddress()
	if err != nil {
		return "", fmt.Errorf("getting chainlink node's primary ETH address failed: %w", err)
	}
	shortNodeAddr := nodeAddress[2:12]
	shortOCRAddr := ocrInstance.Address()[2:12]
	return strings.ToLower(fmt.Sprintf("node_%s_contract_%s", shortNodeAddr, shortOCRAddr)), nil
}
