// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import "../../interfaces/TypeAndVersionInterface.sol";
import "../interfaces/OCR2DRClientInterface.sol";
import "../interfaces/OCR2DROracleInterface.sol";
import "./OCR2DRBillableAbstract.sol";
import "../ocr2/OCR2Base.sol";

/**
 * @title OCR2DR oracle contract
 * @dev THIS CONTRACT HAS NOT GONE THROUGH ANY SECURITY REVIEW. DO NOT USE IN PROD.
 */
contract OCR2DROracle is OCR2DRBillableAbstract, OCR2DROracleInterface, OCR2Base {
  event OracleRequest(bytes32 requestId, bytes data, uint32 gasLimit);
  event OracleResponse(bytes32 requestId);
  event UserCallbackError(bytes32 requestId, string reason);
  event UserCallbackRawError(bytes32 requestId, bytes lowLevelData);

  error EmptyRequestData();
  error InconsistentReportData();
  error EmptyPublicKey();

  bytes private s_donPublicKey;

  constructor() OCR2Base(true) {}

  /**
   * @notice The type and version of this contract
   * @return Type and version string
   */
  function typeAndVersion() external pure override returns (string memory) {
    return "OCR2DROracle 0.0.0";
  }

  /**
   * @inheritdoc OCR2DROracleInterface
   */
  function getRegistry() external view override returns (address) {
    return address(s_registry);
  }

  /**
   * @inheritdoc OCR2DROracleInterface
   */
  function setRegistry(address registryAddress) external override onlyOwner {
    if (registryAddress == address(0)) {
      revert EmptyBillingRegistry();
    }
    s_registry = OCR2DRRegistryInterface(registryAddress);
  }

  /**
   * @inheritdoc OCR2DROracleInterface
   */
  function getDONPublicKey() external view override returns (bytes memory) {
    return s_donPublicKey;
  }

  /**
   * @inheritdoc OCR2DROracleInterface
   */
  function setDONPublicKey(bytes calldata donPublicKey) external override onlyOwner {
    if (donPublicKey.length == 0) {
      revert EmptyPublicKey();
    }
    s_donPublicKey = donPublicKey;
  }

  /**
   * @inheritdoc OCR2DROracleInterface
   */
  function sendRequest(
    uint64 subscriptionId,
    bytes calldata data,
    uint32 gasLimit
  ) external override returns (bytes32) {
    if (data.length == 0) {
      revert EmptyRequestData();
    }
    bytes32 requestId = s_registry.beginBilling(
      data,
      OCR2DRRegistryInterface.RequestBilling(msg.sender, subscriptionId, gasLimit)
    );
    emit OracleRequest(requestId, data, gasLimit);
    return requestId;
  }

  function fulfillRequest(
    bytes32 requestId,
    bytes memory commitment,
    bytes memory response,
    bytes memory err,
    address[maxNumOracles] memory signers,
    uint8 signerCount,
    uint32 reportValidationGas
  ) internal {
    try
      s_registry.concludeBilling(
        requestId,
        commitment,
        response,
        err,
        msg.sender,
        signers,
        signerCount,
        reportValidationGas,
        uint32(gasleft())
      )
    returns (bool success) {
      if (success) {
        emit OracleResponse(requestId);
      } else {
        emit UserCallbackError(requestId, "error in callback");
      }
    } catch (bytes memory reason) {
      emit UserCallbackRawError(requestId, reason);
    }
  }

  function _beforeSetConfig(uint8 _f, bytes memory _onchainConfig) internal override {}

  function _afterSetConfig(uint8 _f, bytes memory _onchainConfig) internal override {}

  function _validateReport(
    bytes32, /* configDigest */
    uint40, /* epochAndRound */
    bytes memory report
  ) internal pure override returns (bool) {
    bytes32[] memory requestIds;
    bytes[] memory results;
    bytes[] memory errors;
    bytes[] memory commitments;
    (requestIds, results, errors, commitments) = abi.decode(report, (bytes32[], bytes[], bytes[], bytes[]));
    if (
      requestIds.length != results.length &&
      requestIds.length != errors.length &&
      requestIds.length != commitments.length
    ) {
      return false;
    }
    return true;
  }

  function _report(
    uint32 initialGas,
    address, /* transmitter */
    uint8 signerCount,
    address[maxNumOracles] memory signers,
    bytes calldata report
  ) internal override {
    bytes32[] memory requestIds;
    bytes[] memory results;
    bytes[] memory errors;
    bytes[] memory commitments;
    (requestIds, results, errors, commitments) = abi.decode(report, (bytes32[], bytes[], bytes[], bytes[]));

    uint256 reportValidationGasShare = (initialGas - gasleft()) / signerCount;

    for (uint256 i = 0; i < requestIds.length; i++) {
      fulfillRequest(
        requestIds[i],
        commitments[i],
        results[i],
        errors[i],
        signers,
        signerCount,
        uint32(reportValidationGasShare)
      );
    }
  }
}
