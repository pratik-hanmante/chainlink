// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import { UpkeepTranscoder4_0, UpkeepV12, UpkeepV13, UpkeepV20} from "../../contracts/UpkeepTranscoder4_0.sol";
import { BaseTest } from "../BaseTest.sol";
import {KeeperRegistryBase2_1 as R21} from "../../KeeperRegistryBase2_1.sol";
import { RegistryVersion, UpkeepV12, UpkeepV13, UpkeepV20 } from "../utils/UpkeepStructs.sol";

contract UpkeepTranscoderSetUp is BaseTest {
    UpkeepTranscoder4_0 internal transcoder;

    function setUp() public {
        transcoder = new UpkeepTranscoder4_0();
    }

    // Add further set-up and test functions here
}

contract UpkeepTranscoder4_0_transcodeUpkeeps is UpkeepTranscoderSetUp {

    // Add your tests here

    function testV12toV21Transcoding() public {
        // Define an example UpkeepV12
        UpkeepV12[] memory upkeepsArrV12;
        upkeepsArrV12[0].balance = 96;
        upkeepsArrV12[0].lastKeeper = address(0);
        upkeepsArrV12[0].executeGas = 32;
        upkeepsArrV12[0].maxValidBlocknumber = 64;
        upkeepsArrV12[0].target = address(0);
        upkeepsArrV12[0].amountSpent = 96;
        upkeepsArrV12[0].admin = address(0);

        // Call the function with the V12 upkeep
        bytes memory result = transcoder.transcodeUpkeeps(uint8(RegistryVersion.V12), uint8(0), abi.encode(upkeepsArrV12));
        
        // Add asserts here to validate the result
        (uint256[] newIds, R21.Upkeep[] memory newUpkeeps, address[] newAdmins, bytes[] memory checkDatas, bytes[] byteArr1, bytes[] byteArr2) = abi.decode(result)
          for (uint256 i = 0; i < ids.length; ix++) {
            assertEq(newIds[i], ids[i])
            // struct equality check
            assertEq(newAdmins[i], upkeepsArrV12[0].admin)
            assertEq()
          }
    }

    // Add similar test functions for V13 and V20
}
