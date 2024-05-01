package measure

import (
	"fmt"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/metrics"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
)

// to test cross-transaction rate
type TestTxNumCount_Relay struct {
	epochID int
	txNum   []float64
}

func NewTestTxNumCount_Relay() *TestTxNumCount_Relay {
	return &TestTxNumCount_Relay{
		epochID: -1,
		txNum:   make([]float64, 0),
	}
}

func (ttnc *TestTxNumCount_Relay) OutputMetricName() string {
	return "Tx_number"
}

func (ttnc *TestTxNumCount_Relay) UpdateMeasureRecord(b *pbft.BlockInfoMsg) {
	if b.BlockBodyLength == 0 { // empty block
		return
	}
	epochid := b.Epoch
	// extend
	for ttnc.epochID < epochid {
		ttnc.txNum = append(ttnc.txNum, 0)
		ttnc.epochID++
	}

	ttnc.txNum[epochid] += float64(len(b.ExcutedTxs))
}

func (ttnc *TestTxNumCount_Relay) OutputRecord() (perEpochCTXs []float64, totTxNum float64) {
	perEpochCTXs = make([]float64, 0)
	totTxNum = 0.0
	for _, tn := range ttnc.txNum {
		perEpochCTXs = append(perEpochCTXs, tn)
		totTxNum += tn
	}
	return perEpochCTXs, totTxNum
}

func (ttnc *TestTxNumCount_Relay) GetDesc() metrics.Desc {
	b := metrics.NewDescBuilder("交易数量统计", "对各个节点交易数量和总交易数量的统计")

	var perEpochCTXs []float64
	var totTxNum float64
	perEpochCTXs = make([]float64, 0)
	totTxNum = 0.0
	for i, tn := range ttnc.txNum {
		b.AddElem(fmt.Sprintf("第%v批次 交易计数（个）", i+1), fmt.Sprintf("第%v批次交易的数量", i+1), tn)
		perEpochCTXs = append(perEpochCTXs, tn)
		totTxNum += tn
	}
	b.AddElem("总交易计数（个）", "总的交易数量统计", totTxNum)
	return b.GetDesc()
}
