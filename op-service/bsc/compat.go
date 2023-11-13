package bsc

import (
	"container/list"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum-optimism/optimism/op-node/eth"
)

var DefaultBaseFee = big.NewInt(3000000000)
var DefaultOPBNBTestnetBaseFee = big.NewInt(5000000000)
var OPBNBTestnet = big.NewInt(5611)
var percentile = 50

const CountBlockSize = 21

type BlockInfoBSCWrapper struct {
	eth.BlockInfo
	baseFee *big.Int
}

var _ eth.BlockInfo = (*BlockInfoBSCWrapper)(nil)

func NewBlockInfoBSCWrapper(info eth.BlockInfo, baseFee *big.Int) *BlockInfoBSCWrapper {
	return &BlockInfoBSCWrapper{
		BlockInfo: info,
		baseFee:   baseFee,
	}
}

func (b *BlockInfoBSCWrapper) BaseFee() *big.Int {
	return b.baseFee
}

// BaseFeeByTransactions calculates the average gas price of the non-zero-gas-price transactions in the block.
// If there is no non-zero-gas-price transaction in the block, it returns DefaultBaseFee.
func BaseFeeByTransactions(transactions types.Transactions) *big.Int {
	nonZeroTxsCnt := big.NewInt(0)
	nonZeroTxsSum := big.NewInt(0)
	for _, tx := range transactions {
		if tx.GasPrice().Cmp(common.Big0) > 0 {
			nonZeroTxsCnt.Add(nonZeroTxsCnt, big.NewInt(1))
			nonZeroTxsSum.Add(nonZeroTxsSum, tx.GasPrice())
		}
	}

	if nonZeroTxsCnt.Cmp(big.NewInt(0)) == 0 {
		return DefaultBaseFee
	}
	return nonZeroTxsSum.Div(nonZeroTxsSum, nonZeroTxsCnt)
}

// BaseFeeByNetworks set l1 gas price by network.
func BaseFeeByNetworks(chainId *big.Int) *big.Int {
	if chainId.Cmp(OPBNBTestnet) == 0 {
		return DefaultOPBNBTestnetBaseFee
	} else {
		return DefaultBaseFee
	}
}

func ToLegacyTx(dynTx *types.DynamicFeeTx) types.TxData {
	return &types.LegacyTx{
		Nonce:    dynTx.Nonce,
		GasPrice: dynTx.GasFeeCap,
		Gas:      dynTx.Gas,
		To:       dynTx.To,
		Value:    dynTx.Value,
		Data:     dynTx.Data,
	}
}

func ToLegacyCallMsg(callMsg ethereum.CallMsg) ethereum.CallMsg {
	return ethereum.CallMsg{
		From:     callMsg.From,
		To:       callMsg.To,
		GasPrice: callMsg.GasFeeCap,
		Data:     callMsg.Data,
	}
}

func MedianGasPrice(transactions types.Transactions) *big.Int {
	var nonZeroTxsGasPrice []*big.Int
	for _, tx := range transactions {
		if tx.GasPrice().Cmp(common.Big0) > 0 {
			nonZeroTxsGasPrice = append(nonZeroTxsGasPrice, tx.GasPrice())
		}
	}
	sort.Sort(bigIntArray(nonZeroTxsGasPrice))
	medianGasPrice := DefaultBaseFee
	if len(nonZeroTxsGasPrice) != 0 {
		medianGasPrice = nonZeroTxsGasPrice[(len(nonZeroTxsGasPrice)-1)*percentile/100]
	}
	return medianGasPrice
}

func FinalGasPrice(medianGasPriceQueue list.List) *big.Int {
	var allMedianGasPrice []*big.Int
	for item := medianGasPriceQueue.Front(); item != nil; item = item.Next() {
		allMedianGasPrice = append(allMedianGasPrice, item.Value.(*big.Int))
	}
	sort.Sort(bigIntArray(allMedianGasPrice))
	finalGasPrice := allMedianGasPrice[(len(allMedianGasPrice)-1)*percentile/100]
	medianGasPriceQueue.Remove(medianGasPriceQueue.Front())
	return finalGasPrice
}

type bigIntArray []*big.Int

func (s bigIntArray) Len() int           { return len(s) }
func (s bigIntArray) Less(i, j int) bool { return s[i].Cmp(s[j]) < 0 }
func (s bigIntArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
