pragma solidity 0.8.6;

import "../BaseTest.t.sol";
import {VRF} from "../../../../src/v0.8/vrf/VRF.sol";
import {MockLinkToken} from "../../../../src/v0.8/mocks/MockLinkToken.sol";
import {MockV3Aggregator} from "../../../../src/v0.8/tests/MockV3Aggregator.sol";
import {ExposedVRFCoordinatorV2Plus} from "../../../../src/v0.8/dev/vrf/testhelpers/ExposedVRFCoordinatorV2Plus.sol";
import {VRFCoordinatorV2Plus} from "../../../../src/v0.8/dev/vrf/VRFCoordinatorV2Plus.sol";
import {BlockhashStore} from "../../../../src/v0.8/dev/BlockhashStore.sol";
import {VRFV2PlusConsumerExample} from "../../../../src/v0.8/dev/vrf/testhelpers/VRFV2PlusConsumerExample.sol";
import {console} from "forge-std/console.sol";

/*
 * USAGE INSTRUCTIONS:
 * To add new tests/proofs, uncomment the "console.sol" import from foundry, and gather key fields
 * from your VRF request.
 * Then, pass your request info into the generate-proof-v2-plus script command
 * located in /core/scripts/vrfv2/testnet/proofs.go to generate a proof that can be tested on-chain.
 **/

contract VRFV2Plus is BaseTest {
  address internal constant LINK_WHALE = 0xD883a6A1C22fC4AbFE938a5aDF9B2Cc31b1BF18B;

  // Bytecode for a VRFV2PlusConsumerExample contract.
  // to calculate: console.logBytes(type(VRFV2PlusConsumerExample).creationCode);
  bytes constant initializeCode =
    hex"60806040523480156200001157600080fd5b50604051620014ea380380620014ea833981016040819052620000349162000213565b8133806000816200008c5760405162461bcd60e51b815260206004820152601860248201527f43616e6e6f7420736574206f776e657220746f207a65726f000000000000000060448201526064015b60405180910390fd5b600080546001600160a01b0319166001600160a01b0384811691909117909155811615620000bf57620000bf816200014a565b5050506001600160a01b038116620001095760405162461bcd60e51b815260206004820152600c60248201526b7a65726f206164647265737360a01b604482015260640162000083565b600280546001600160a01b03199081166001600160a01b0393841617909155600480548216948316949094179093556005805490931691161790556200024b565b6001600160a01b038116331415620001a55760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640162000083565b600180546001600160a01b0319166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b80516001600160a01b03811681146200020e57600080fd5b919050565b600080604083850312156200022757600080fd5b6200023283620001f6565b91506200024260208401620001f6565b90509250929050565b61128f806200025b6000396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c806379ba509711610097578063a168fa8911610066578063a168fa8914610220578063cf62c8ab14610280578063eff2701714610293578063f2fde38b146102a657600080fd5b806379ba5097146101e15780637ec8773a146101e95780638da5cb5b146101fc5780639eccacf61461020d57600080fd5b806344ff81ce116100d357806344ff81ce146101665780635d7d53e314610179578063706da1ca146101825780637725135b146101b657600080fd5b80631fe543e31461010557806329e5d8311461011a5780632fa4e4421461014057806336bfffed14610153575b600080fd5b610118610113366004610f46565b6102b9565b005b61012d610128366004610fea565b610325565b6040519081526020015b60405180910390f35b61011861014e3660046110a1565b61043b565b610118610161366004610e53565b6104f1565b610118610174366004610e31565b610623565b61012d60065481565b60055461019d90600160a01b900467ffffffffffffffff1681565b60405167ffffffffffffffff9091168152602001610137565b6005546101c9906001600160a01b031681565b6040516001600160a01b039091168152602001610137565b6101186106ab565b6101186101f7366004610e31565b610769565b6000546001600160a01b03166101c9565b6004546101c9906001600160a01b031681565b61025b61022e366004610f14565b6007602052600090815260409020805460019091015460ff82169161010090046001600160a01b03169083565b6040805193151584526001600160a01b03909216602084015290820152606001610137565b61011861028e3660046110a1565b610775565b6101186102a136600461100c565b6108ef565b6101186102b4366004610e31565b610add565b6002546001600160a01b03163314610317576002546040517f1cf993f40000000000000000000000000000000000000000000000000000000081523360048201526001600160a01b0390911660248201526044015b60405180910390fd5b6103218282610aee565b5050565b60008281526007602090815260408083208151608081018352815460ff81161515825261010090046001600160a01b031681850152600182015481840152600282018054845181870281018701909552808552869592946060860193909291908301828280156103b457602002820191906000526020600020905b8154815260200190600101908083116103a0575b50505050508152505090508060400151600014156104145760405162461bcd60e51b815260206004820152601760248201527f7265717565737420494420697320696e636f7272656374000000000000000000604482015260640161030e565b8060600151838151811061042a5761042a611248565b602002602001015191505092915050565b60055460045460408051600160a01b840467ffffffffffffffff1660208201526001600160a01b0393841693634000aea09316918591016040516020818303038152906040526040518463ffffffff1660e01b815260040161049f9392919061111c565b602060405180830381600087803b1580156104b957600080fd5b505af11580156104cd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103219190610ef7565b600554600160a01b900467ffffffffffffffff166105515760405162461bcd60e51b815260206004820152600d60248201527f7375624944206e6f742073657400000000000000000000000000000000000000604482015260640161030e565b60005b81518110156103215760045460055483516001600160a01b0390921691637341c10c91600160a01b900467ffffffffffffffff169085908590811061059b5761059b611248565b60200260200101516040518363ffffffff1660e01b81526004016105de92919067ffffffffffffffff9290921682526001600160a01b0316602082015260400190565b600060405180830381600087803b1580156105f857600080fd5b505af115801561060c573d6000803e3d6000fd5b50505050808061061b9061121f565b915050610554565b6003546001600160a01b0316331461067c576003546040517f4ae338ff0000000000000000000000000000000000000000000000000000000081523360048201526001600160a01b03909116602482015260440161030e565b6002805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055565b6001546001600160a01b031633146107055760405162461bcd60e51b815260206004820152601660248201527f4d7573742062652070726f706f736564206f776e657200000000000000000000604482015260640161030e565b600080543373ffffffffffffffffffffffffffffffffffffffff19808316821784556001805490911690556040516001600160a01b0390921692909183917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a350565b61077281610b81565b50565b600554600160a01b900467ffffffffffffffff1661043b5760048054604080517fa21a23e400000000000000000000000000000000000000000000000000000000815290516001600160a01b039092169263a21a23e49282820192602092908290030181600087803b1580156107ea57600080fd5b505af11580156107fe573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108229190611077565b600580547fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff16600160a01b67ffffffffffffffff93841681029190911791829055600480546040517f7341c10c00000000000000000000000000000000000000000000000000000000815292909304909316928101929092523060248301526001600160a01b031690637341c10c90604401600060405180830381600087803b1580156108ce57600080fd5b505af11580156108e2573d6000803e3d6000fd5b5050505061043b30610b81565b60006040518060c00160405280848152602001600560149054906101000a900467ffffffffffffffff1667ffffffffffffffff1681526020018661ffff1681526020018763ffffffff1681526020018563ffffffff1681526020016109636040518060200160405280861515815250610bf0565b9052600480546040517f596b8b880000000000000000000000000000000000000000000000000000000081529293506000926001600160a01b039091169163596b8b88916109b39186910161115b565b602060405180830381600087803b1580156109cd57600080fd5b505af11580156109e1573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a059190610f2d565b604080516080810182526000808252336020808401918252838501868152855184815280830187526060860190815287855260078352959093208451815493517fffffffffffffffffffffff0000000000000000000000000000000000000000009094169015157fffffffffffffffffffffff0000000000000000000000000000000000000000ff16176101006001600160a01b039094169390930292909217825591516001820155925180519495509193849392610acb926002850192910190610da1565b50505060069190915550505050505050565b610ae5610c8e565b61077281610cea565b6006548214610b3f5760405162461bcd60e51b815260206004820152601760248201527f7265717565737420494420697320696e636f7272656374000000000000000000604482015260640161030e565b60008281526007602090815260409091208251610b6492600290920191840190610da1565b50506000908152600760205260409020805460ff19166001179055565b6001600160a01b038116610bc1576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6003805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055565b60607f92fd13387c7fe7befbc38d303d6468778fb9731bc4583f17d92989c6fcfdeaaa82604051602401610c2991511515815260200190565b60408051601f198184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff000000000000000000000000000000000000000000000000000000009093169290921790915292915050565b6000546001600160a01b03163314610ce85760405162461bcd60e51b815260206004820152601660248201527f4f6e6c792063616c6c61626c65206279206f776e657200000000000000000000604482015260640161030e565b565b6001600160a01b038116331415610d435760405162461bcd60e51b815260206004820152601760248201527f43616e6e6f74207472616e7366657220746f2073656c66000000000000000000604482015260640161030e565b6001805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0383811691821790925560008054604051929316917fed8889f560326eb138920d842192f0eb3dd22b4f139c87a2c57538e05bae12789190a350565b828054828255906000526020600020908101928215610ddc579160200282015b82811115610ddc578251825591602001919060010190610dc1565b50610de8929150610dec565b5090565b5b80821115610de85760008155600101610ded565b80356001600160a01b0381168114610e1857600080fd5b919050565b803563ffffffff81168114610e1857600080fd5b600060208284031215610e4357600080fd5b610e4c82610e01565b9392505050565b60006020808385031215610e6657600080fd5b823567ffffffffffffffff811115610e7d57600080fd5b8301601f81018513610e8e57600080fd5b8035610ea1610e9c826111fb565b6111ca565b80828252848201915084840188868560051b8701011115610ec157600080fd5b600094505b83851015610eeb57610ed781610e01565b835260019490940193918501918501610ec6565b50979650505050505050565b600060208284031215610f0957600080fd5b8151610e4c81611274565b600060208284031215610f2657600080fd5b5035919050565b600060208284031215610f3f57600080fd5b5051919050565b60008060408385031215610f5957600080fd5b8235915060208084013567ffffffffffffffff811115610f7857600080fd5b8401601f81018613610f8957600080fd5b8035610f97610e9c826111fb565b80828252848201915084840189868560051b8701011115610fb757600080fd5b600094505b83851015610fda578035835260019490940193918501918501610fbc565b5080955050505050509250929050565b60008060408385031215610ffd57600080fd5b50508035926020909101359150565b600080600080600060a0868803121561102457600080fd5b61102d86610e1d565b9450602086013561ffff8116811461104457600080fd5b935061105260408701610e1d565b925060608601359150608086013561106981611274565b809150509295509295909350565b60006020828403121561108957600080fd5b815167ffffffffffffffff81168114610e4c57600080fd5b6000602082840312156110b357600080fd5b81356bffffffffffffffffffffffff81168114610e4c57600080fd5b6000815180845260005b818110156110f5576020818501810151868301820152016110d9565b81811115611107576000602083870101525b50601f01601f19169290920160200192915050565b6001600160a01b03841681526bffffffffffffffffffffffff8316602082015260606040820152600061115260608301846110cf565b95945050505050565b602081528151602082015267ffffffffffffffff602083015116604082015261ffff60408301511660608201526000606083015163ffffffff80821660808501528060808601511660a0850152505060a083015160c0808401526111c260e08401826110cf565b949350505050565b604051601f8201601f1916810167ffffffffffffffff811182821017156111f3576111f361125e565b604052919050565b600067ffffffffffffffff8211156112155761121561125e565b5060051b60200190565b600060001982141561124157634e487b7160e01b600052601160045260246000fd5b5060010190565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052604160045260246000fd5b801515811461077257600080fdfea164736f6c6343000806000a";

  BlockhashStore s_bhs;
  ExposedVRFCoordinatorV2Plus s_testCoordinator;
  VRFV2PlusConsumerExample s_testConsumer;
  MockLinkToken s_linkToken;
  MockV3Aggregator s_linkEthFeed;

  VRFCoordinatorV2Plus.FeeConfig basicFeeConfig =
    VRFCoordinatorV2Plus.FeeConfig({fulfillmentFlatFeeLinkPPM: 0, fulfillmentFlatFeeEthPPM: 0});

  // VRF KeyV2 generated from a node; not sensitive information.
  // The secret key used to generate this key is: 10.
  bytes vrfUncompressedPublicKey =
    hex"a0434d9e47f3c86235477c7b1ae6ae5d3442d49b1943c2b752a68e2a47e247c7893aba425419bc27a3b6c7e693a24c696f794c2ed877a1593cbee53b037368d7";
  bytes vrfCompressedPublicKey = hex"a0434d9e47f3c86235477c7b1ae6ae5d3442d49b1943c2b752a68e2a47e247c701";
  bytes32 vrfKeyHash = hex"9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528";

  function setUp() public override {
    BaseTest.setUp();

    // Fund our users.
    vm.roll(1);
    vm.deal(LINK_WHALE, 10_000 ether);
    changePrank(LINK_WHALE);

    // Instantiate BHS.
    s_bhs = new BlockhashStore();

    // Deploy coordinator and consumer.
    s_testCoordinator = new ExposedVRFCoordinatorV2Plus(address(s_bhs));

    // Deploy link token and link/eth feed.
    s_linkToken = new MockLinkToken();
    s_linkEthFeed = new MockV3Aggregator(18, 500000000000000000); // .5 ETH (good for testing)

    // Use create2 to deploy our consumer, so that its address is always the same
    // and surrounding changes do not alter our generated proofs.
    bytes memory consumerInitCode = bytes.concat(
      initializeCode,
      abi.encode(address(s_testCoordinator), address(s_linkToken))
    );
    bytes32 abiEncodedOwnerAddress = bytes32(uint256(uint160(LINK_WHALE)) << 96);
    address consumerCreate2Address;
    assembly {
      consumerCreate2Address := create2(
        0, // value - left at zero here
        add(0x20, consumerInitCode), // initialization bytecode (excluding first memory slot which contains its length)
        mload(consumerInitCode), // length of initialization bytecode
        abiEncodedOwnerAddress // user-defined nonce to ensure unique SCA addresses
      )
    }
    s_testConsumer = VRFV2PlusConsumerExample(consumerCreate2Address);

    // Configure the coordinator.
    s_testCoordinator.setLINK(address(s_linkToken));
    s_testCoordinator.setLinkEthFeed(address(s_linkEthFeed));
  }

  function setConfig(VRFCoordinatorV2Plus.FeeConfig memory feeConfig) internal {
    s_testCoordinator.setConfig(
      0, // minRequestConfirmations
      2_500_000, // maxGasLimit
      1, // stalenessSeconds
      50_000, // gasAfterPaymentCalculation
      50000000000000000, // fallbackWeiPerUnitLink
      feeConfig
    );
  }

  function testSetConfig() public {
    // Should setConfig successfully.
    setConfig(basicFeeConfig);
    (uint16 minConfs, uint32 gasLimit, ) = s_testCoordinator.getRequestConfig();
    assertEq(minConfs, 0);
    assertEq(gasLimit, 2_500_000);

    // Test that setting requestConfirmations above MAX_REQUEST_CONFIRMATIONS reverts.
    vm.expectRevert(abi.encodeWithSelector(VRFCoordinatorV2Plus.InvalidRequestConfirmations.selector, 500, 500, 200));
    s_testCoordinator.setConfig(500, 2_500_000, 1, 50_000, 50000000000000000, basicFeeConfig);

    // Test that setting fallbackWeiPerUnitLink to zero reverts.
    vm.expectRevert(abi.encodeWithSelector(VRFCoordinatorV2Plus.InvalidLinkWeiPrice.selector, 0));
    s_testCoordinator.setConfig(0, 2_500_000, 1, 50_000, 0, basicFeeConfig);
  }

  function testRegisterProvingKey() public {
    // Should set the proving key successfully.
    registerProvingKey();
    (, , bytes32[] memory keyHashes) = s_testCoordinator.getRequestConfig();
    assertEq(keyHashes[0], vrfKeyHash);

    // Should revert when already registered.
    uint256[2] memory uncompressedKeyParts = this.getProvingKeyParts(vrfUncompressedPublicKey);
    vm.expectRevert(abi.encodeWithSelector(VRFCoordinatorV2Plus.ProvingKeyAlreadyRegistered.selector, vrfKeyHash));
    s_testCoordinator.registerProvingKey(LINK_WHALE, uncompressedKeyParts);
  }

  function registerProvingKey() public {
    uint256[2] memory uncompressedKeyParts = this.getProvingKeyParts(vrfUncompressedPublicKey);
    s_testCoordinator.registerProvingKey(LINK_WHALE, uncompressedKeyParts);
  }

  // note: Call this function via this.getProvingKeyParts to be able to pass memory as calldata and
  // index over the byte array.
  function getProvingKeyParts(bytes calldata uncompressedKey) public pure returns (uint256[2] memory) {
    uint256 keyPart1 = uint256(bytes32(uncompressedKey[0:32]));
    uint256 keyPart2 = uint256(bytes32(uncompressedKey[32:64]));
    return [keyPart1, keyPart2];
  }

  function testCreateSubscription() public {
    uint64 subId = s_testCoordinator.createSubscription();
    s_testCoordinator.fundSubscriptionWithEth{value: 10 ether}(subId);
  }

  event RandomWordsRequested(
    bytes32 indexed keyHash,
    uint256 requestId,
    uint256 preSeed,
    uint64 indexed subId,
    uint16 minimumRequestConfirmations,
    uint32 callbackGasLimit,
    uint32 numWords,
    bool nativePayment,
    address indexed sender
  );

  function testRequestAndFulfillRandomWordsNative() public {
    uint32 requestBlock = 10;
    vm.roll(requestBlock);
    s_testConsumer.createSubscriptionAndFund(0);
    s_testConsumer.setSubOwner(LINK_WHALE);
    uint64 subId = s_testConsumer.s_subId();
    s_testCoordinator.fundSubscriptionWithEth{value: 10 ether}(subId);

    // Apply basic configs to contract.
    setConfig(basicFeeConfig);
    registerProvingKey();

    // Request random words.
    vm.expectEmit(true, true, true, true);
    (uint256 requestId, uint256 preSeed) = s_testCoordinator.computeRequestIdExternal(
      vrfKeyHash,
      address(s_testConsumer),
      subId,
      2
    );
    emit RandomWordsRequested(
      vrfKeyHash,
      requestId,
      preSeed,
      1, // subId
      0, // minConfirmations
      1_000_000, // callbackGasLimit
      1, // numWords
      true, // nativePayment
      address(s_testConsumer) // requester
    );
    s_testConsumer.requestRandomWords(1_000_000, 0, 1, vrfKeyHash, true);
    (bool fulfilled, , ) = s_testConsumer.s_requests(requestId);
    assertEq(fulfilled, false);

    // Uncomment these console logs to see info about the request:
    // console.log("requestId: ", requestId);
    // console.log("preSeed: ", preSeed);
    // console.log("sender: ", address(s_testConsumer));

    // Move on to the next block.
    // Store the previous block's blockhash, and assert that it is as expected.
    vm.roll(requestBlock + 1);
    s_bhs.store(requestBlock);
    assertEq(hex"c65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a8", s_bhs.getBlockhash(requestBlock));

    // Fulfill the request.
    // Proof generated via the generate-proof-v2-plus script command. Example usage:
    /*
        go run . generate-proof-v2-plus \
        -key-hash 0x9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528 \
        -pre-seed 75691120191018016735882771624817645647824742686222051289703807234819010996076 \
        -block-hash 0xc65a7bb8d6351c1cf70c95a316cc6a92839c986682d98bc35f958f4883f9d2a8 \
        -block-num 10 \
        -sender 0xc8E6EbCA64Bfb55c629377B6A13960A81f40B581 \
        -native-payment true
        */
    VRF.Proof memory proof = VRF.Proof({
      pk: [
        72488970228380509287422715226575535698893157273063074627791787432852706183111,
        62070622898698443831883535403436258712770888294397026493185421712108624767191
      ],
      gamma: [
        75383267251969672296268498817560629288818643643730140607865534793758789115008,
        56317310071450638585325149087394219231772751102251183900642124310206345507576
      ],
      c: 61647911377666443996029446682766202269810624765940572348506681098037871800476,
      s: 100567053859228822319671051136344708045292695533678153813883547397557870986442,
      seed: 75691120191018016735882771624817645647824742686222051289703807234819010996076,
      uWitness: 0xDDcC42C282aa7B67D7e2d09C2e98d2B1042123Fd,
      cGammaWitness: [
        79158763384243139474431077481157249695173423054401214119485945207803480097214,
        100312278295182067690393476082104650895185768438478876495894908427697028038032
      ],
      sHashWitness: [
        44558937922788029349366556709547001403757954392754660222266934662747487173199,
        51903593851711669113005633123615189441216700946808828391159740573376725220843
      ],
      zInv: 80260799599166078749892682204088116396346875228077623955725975464581138368082
    });
    VRFCoordinatorV2Plus.RequestCommitment memory rc = VRFCoordinatorV2Plus.RequestCommitment({
      blockNum: requestBlock,
      subId: 1,
      callbackGasLimit: 1_000_000,
      numWords: 1,
      sender: address(s_testConsumer),
      nativePayment: true
    });
    (, uint96 ethBalanceBefore, , ) = s_testCoordinator.getSubscription(subId);
    s_testCoordinator.fulfillRandomWords{gas: 1_500_000}(proof, rc);
    (fulfilled, , ) = s_testConsumer.s_requests(requestId);
    assertEq(fulfilled, true);

    // The cost of fulfillRandomWords is approximately 100_000 gas.
    // gasAfterPaymentCalculation is 50_000.
    //
    // The cost of the VRF fulfillment charged to the user is:
    // baseFeeWei = weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft())
    // baseFeeWei = 1 * (50_000 + 100_000)
    // baseFeeWei = 150_000
    // ...
    // billed_fee = baseFeeWei + flatFeeWei + l1CostWei
    // billed_fee = baseFeeWei + 0 + 0
    // billed_fee = 150_000
    (, uint96 ethBalanceAfter, , ) = s_testCoordinator.getSubscription(subId);
    assertApproxEqAbs(ethBalanceAfter, ethBalanceBefore - 120_000, 10_000);
  }

  function testRequestAndFulfillRandomWordsLINK() public {
    uint32 requestBlock = 20;
    vm.roll(requestBlock);
    s_linkToken.transfer(address(s_testConsumer), 10 ether);
    s_testConsumer.createSubscriptionAndFund(10 ether);
    s_testConsumer.setSubOwner(LINK_WHALE);
    uint64 subId = s_testConsumer.s_subId();

    // Apply basic configs to contract.
    setConfig(basicFeeConfig);
    registerProvingKey();

    // Request random words.
    vm.expectEmit(true, true, true, true);
    (uint256 requestId, uint256 preSeed) = s_testCoordinator.computeRequestIdExternal(
      vrfKeyHash,
      address(s_testConsumer),
      subId,
      2
    );
    emit RandomWordsRequested(
      vrfKeyHash,
      requestId,
      preSeed,
      1, // subId
      0, // minConfirmations
      1_000_000, // callbackGasLimit
      1, // numWords
      false, // nativePayment
      address(s_testConsumer) // requester
    );
    s_testConsumer.requestRandomWords(1_000_000, 0, 1, vrfKeyHash, false);
    (bool fulfilled, , ) = s_testConsumer.s_requests(requestId);
    assertEq(fulfilled, false);

    // Uncomment these console logs to see info about the request:
    // console.log("requestId: ", requestId);
    // console.log("preSeed: ", preSeed);
    // console.log("sender: ", address(s_testConsumer));

    // Move on to the next block.
    // Store the previous block's blockhash, and assert that it is as expected.
    vm.roll(requestBlock + 1);
    s_bhs.store(requestBlock);
    assertEq(hex"ce6d7b5282bd9a3661ae061feed1dbda4e52ab073b1f9285be6e155d9c38d4ec", s_bhs.getBlockhash(requestBlock));

    // Fulfill the request.
    // Proof generated via the generate-proof-v2-plus script command. Example usage:
    /*
        go run . generate-proof-v2-plus \
        -key-hash 0x9f2353bde94264dbc3d554a94cceba2d7d2b4fdce4304d3e09a1fea9fbeb1528 \
        -pre-seed 75691120191018016735882771624817645647824742686222051289703807234819010996076 \
        -block-hash 0xce6d7b5282bd9a3661ae061feed1dbda4e52ab073b1f9285be6e155d9c38d4ec \
        -block-num 20 \
        -sender 0xc8E6EbCA64Bfb55c629377B6A13960A81f40B581
    */
    VRF.Proof memory proof = VRF.Proof({
      pk: [
        72488970228380509287422715226575535698893157273063074627791787432852706183111,
        62070622898698443831883535403436258712770888294397026493185421712108624767191
      ],
      gamma: [
        31461201777511708081918305128835245000666140073976840051479388569618603648903,
        68844459275635126487349613957273843734172660932713066865162200984073737886836
      ],
      c: 4007610028535731166811646741680870477426492196518841063133182608608788308819,
      s: 96098842852617677101456878862714476238560342792960218463862001840420419346845,
      seed: 75691120191018016735882771624817645647824742686222051289703807234819010996076,
      uWitness: 0x11962D078fABDeB8864165E430476200953f4c6E,
      cGammaWitness: [
        77521242639329621821630885273458534406377967750646690452497670019893779096803,
        62348556077821755766997765471730575865600301131535212595841526930927520869016
      ],
      sHashWitness: [
        51796806133116806660184658627217367381780582813423685829170662901659971854870,
        6387413072742906621199621824698662262666827884361549573353608319121858479482
      ],
      zInv: 55470199074701968783284118261557709031321707846498128347209444775968786079604
    });
    VRFCoordinatorV2Plus.RequestCommitment memory rc = VRFCoordinatorV2Plus.RequestCommitment({
      blockNum: requestBlock,
      subId: 1,
      callbackGasLimit: 1000000,
      numWords: 1,
      sender: address(s_testConsumer),
      nativePayment: false
    });
    (uint96 linkBalanceBefore, , , ) = s_testCoordinator.getSubscription(subId);
    s_testCoordinator.fulfillRandomWords{gas: 1_500_000}(proof, rc);
    (fulfilled, , ) = s_testConsumer.s_requests(requestId);
    assertEq(fulfilled, true);

    // The cost of fulfillRandomWords is approximately 90_000 gas.
    // gasAfterPaymentCalculation is 50_000.
    //
    // The cost of the VRF fulfillment charged to the user is:
    // paymentNoFee = (weiPerUnitGas * (gasAfterPaymentCalculation + startGas - gasleft() + l1CostWei) / link_eth_ratio)
    // paymentNoFee = (1 * (50_000 + 90_000 + 0)) / .5
    // paymentNoFee = 280_000
    // ...
    // billed_fee = paymentNoFee + fulfillmentFlatFeeLinkPPM
    // billed_fee = baseFeeWei + 0
    // billed_fee = 280_000
    // note: delta is doubled from the native test to account for more variance due to the link/eth ratio
    (uint96 linkBalanceAfter, , , ) = s_testCoordinator.getSubscription(subId);
    assertApproxEqAbs(linkBalanceAfter, linkBalanceBefore - 280_000, 20_000);
  }
}
