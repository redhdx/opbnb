// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import { CommonTest } from "./CommonTest.t.sol";
import { L2Config } from "../L2/L2Config.sol";
import { Constants } from "../libraries/Constants.sol";

contract L2Config_Init is CommonTest {
    L2Config l2Config;

    function setUp() public virtual override {
        super.setUp();

        l2Config = new L2Config({
            _owner: alice,
            _l1MinGasPrice: 3 gwei
        });
    }
}

contract L2Config_Setters_TestFail is L2Config_Init {
    function test_setL1MinGasPrice_notOwner_reverts() external {
        vm.expectRevert("Ownable: caller is not the owner");
        l2Config.setL1MinGasPrice(0);
    }
}

contract L2Config_Setters_Test is L2Config_Init {
    event ConfigUpdate(
        uint256 indexed version,
        L2Config.UpdateType indexed updateType,
        bytes data
    );

    function testFuzz_setGasLimit_succeeds(uint64 newL1MinGasPrice) external {
        vm.expectEmit(true, true, true, true);
        emit ConfigUpdate(0, L2Config.UpdateType.L1MINGASPRICE, abi.encode(newL1MinGasPrice));

        vm.prank(l2Config.owner());
        l2Config.setL1MinGasPrice(newL1MinGasPrice);
        assertEq(l2Config.l1MinGasPrice(), newL1MinGasPrice);
    }
}
