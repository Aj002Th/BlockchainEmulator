package measure

import (
	"fmt"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/metrics"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
)

type CrossTxRate struct {
	epochID       int
	totTxNum      []float64
	totCrossTxNum []float64
}

func NewTestCrossTxRate_Pbft() *CrossTxRate {
	return &CrossTxRate{
		epochID:       -1,
		totTxNum:      make([]float64, 0),
		totCrossTxNum: make([]float64, 0),
	}
}

func (tctr *CrossTxRate) OutputMetricName() string {
	return "CrossTransaction_ratio"
}

func (tctr *CrossTxRate) UpdateMeasureRecord(b *pbft.BlockInfoMsg) {
	if b.BlockBodyLength == 0 {
		return
	}
	epochID := b.Epoch

	for tctr.epochID < epochID {
		tctr.totTxNum = append(tctr.totTxNum, 0)
		tctr.totCrossTxNum = append(tctr.totCrossTxNum, 0)
		tctr.epochID++
	}

	for _, _ = range b.ExcutedTxs {
		tctr.totTxNum[epochID] += 0.5
		tctr.totCrossTxNum[epochID] += 0.5

	}
}

func (tctr *CrossTxRate) OutputRecord() (perEpochCTXratio []float64, totCTXratio float64) {
	perEpochCTXratio = make([]float64, 0)
	allEpoch_totTxNum := 0.0
	allEpoch_ctxNum := 0.0
	for eid, totTxN := range tctr.totTxNum {
		perEpochCTXratio = append(perEpochCTXratio, tctr.totCrossTxNum[eid]/totTxN)
		allEpoch_totTxNum += totTxN
		allEpoch_ctxNum += tctr.totCrossTxNum[eid]
	}
	perEpochCTXratio = append(perEpochCTXratio, allEpoch_totTxNum)
	perEpochCTXratio = append(perEpochCTXratio, allEpoch_ctxNum)

	return perEpochCTXratio, allEpoch_ctxNum / allEpoch_totTxNum
}

func (tctr *CrossTxRate) GetDesc() metrics.Desc {
	b := metrics.NewDescBuilder("跨交易率(CrossTxRate)", "平均每秒产生的交易，衡量交易的次数。单位为 交易/秒")

	var perEpochCTXratio []float64

	perEpochCTXratio = make([]float64, 0)
	allEpochTotTxNum := 0.0
	allEpochCtxNum := 0.0
	for eid, totTxN := range tctr.totTxNum {
		b.AddElem(fmt.Sprintf("第%v批次 跨交易率", eid), "", tctr.totCrossTxNum[eid]/totTxN)
		perEpochCTXratio = append(perEpochCTXratio, tctr.totCrossTxNum[eid]/totTxN)
		allEpochTotTxNum += totTxN
		allEpochCtxNum += tctr.totCrossTxNum[eid]
	}
	perEpochCTXratio = append(perEpochCTXratio, allEpochTotTxNum)
	perEpochCTXratio = append(perEpochCTXratio, allEpochCtxNum)
	b.AddElem("交易数量总计", "All Epoch Total Tx Num", allEpochTotTxNum)
	b.AddElem("跨交易数量总计", "", allEpochCtxNum)

	b.AddElem("总体跨交易率", "整个运行过程中的跨交易比率", allEpochCtxNum/allEpochTotTxNum)
	return b.GetDesc()
}
