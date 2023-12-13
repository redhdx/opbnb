// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import {
    OwnableUpgradeable
} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import { Semver } from "../universal/Semver.sol";

/**
 * @title L2Config
 * @notice The L2Config contract is used to manage configuration of opBNB network.
 */
contract L2Config is OwnableUpgradeable, Semver {
    /**
     * @notice Enum representing different types of updates.
     *
     * @custom:value L1minGasPrice  Represents an update to the l1 min gas price.
     */
    enum UpdateType {
        L1MINGASPRICE
    }

    /**
 * @notice Version identifier, used for upgrades.
     */
    uint256 public constant VERSION = 0;

    /**
     * @notice L1 min gas price.
     */
    uint256 public l1MinGasPrice;

    /**
     * @notice Emitted when configuration is updated
     *
     * @param version    L2Config version.
     * @param updateType Type of update.
     * @param data       Encoded update data.
     */
    event ConfigUpdate(uint256 indexed version, UpdateType indexed updateType, bytes data);

    /**
     * @custom:semver 1.0.0
     *
     * @param _owner             Initial owner of the contract.
     * @param _l1MinGasPrice     Initial l1 min gas price.
     */
    constructor(
        address _owner,
        uint64 _l1MinGasPrice
    ) Semver(1, 0, 0) {
        initialize({
            _owner: _owner,
            _l1MinGasPrice: _l1MinGasPrice
        });
    }

    /**
     * @notice Initializer.
     *
     * @param _owner             Initial owner of the contract.
     * @param _l1MinGasPrice     Initial l1 min gas price.
     */
    function initialize(
        address _owner,
        uint64 _l1MinGasPrice
    ) public initializer {
        __Ownable_init();
        transferOwnership(_owner);
        l1MinGasPrice = _l1MinGasPrice;
    }

    /**
     * @notice Updates the L1 min gas price.
     *
     * @param _l1MinGasPrice New L1 min gas price.
     */
    function setL1MinGasPrice(uint64 _l1MinGasPrice) external onlyOwner {
        l1MinGasPrice = _l1MinGasPrice;
        bytes memory data = abi.encode(_l1MinGasPrice);
        emit ConfigUpdate(VERSION, UpdateType.L1MINGASPRICE, data);
    }
}
