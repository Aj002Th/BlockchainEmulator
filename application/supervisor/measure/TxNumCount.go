package measure

import (
	"fmt"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/metrics"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
)

// to test cross-transaction rate
type TxNumCount struct {
	epochID int
	txNum   []float64
}

func NewTxNumCount() *TxNumCount {
	return &TxNumCount{
		epochID: -1,
		txNum:   make([]float64, 0),
	}
}

func (ttnc *TxNumCount) OutputMetricName() string {
	return "Tx_number"
}

func (ttnc *TxNumCount) UpdateMeasureRecord(b *pbft.BlockInfoMsg) {
	if b.BlockBodyLength == 0 { // empty block
		return
	}
	epochID := b.Epoch
	// extend
	for ttnc.epochID < epochID {
		ttnc.txNum = append(ttnc.txNum, 0)
		ttnc.epochID++
	}

	ttnc.txNum[epochID] += float64(len(b.ExcutedTxs))
}

func (ttnc *TxNumCount) OutputRecord() (perEpochCTXs []float64, totTxNum float64) {
	perEpochCTXs = make([]float64, 0)
	totTxNum = 0.0
	for _, tn := range ttnc.txNum {
		perEpochCTXs = append(perEpochCTXs, tn)
		totTxNum += tn
	}
	return perEpochCTXs, totTxNum
}

func (ttnc *TxNumCount) GetDesc() metrics.Desc {
	b := metrics.NewDescBuilder("交易数量统计", "对各个节点交易数量和总交易数量的统计")

	var perEpochCTXs []float64
	var totTxNum float64
	perEpochCTXs = make([]float64, 0)
	totTxNum = 0.0
	for i, tn := range ttnc.txNum {
		b.AddElem(fmt.Sprintf("第%v批次 交易计数", i+1), fmt.Sprintf("第%v批次交易的数量", i+1), fmt.Sprintf("%v 个", tn))
		perEpochCTXs = append(perEpochCTXs, tn)
		totTxNum += tn
	}
	b.AddElem("总交易计数", "总的交易数量统计", fmt.Sprintf("%v 个", totTxNum))
	return b.GetDesc()
}
