package measure

import (
	"fmt"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/metrics"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
)

// to test average Transaction_Confirm_Latency (TCL)  in this system
type TestModule_TCL_Pbft struct {
	epochID           int       // 其实这应该叫epochCnt，随着bim到来，这个值是计数递增的。
	totTxLatencyEpoch []float64 // record the Transaction_Confirm_Latency in each epoch
	txNum             []float64 // record the txNumber in each epoch
}

func NewTestModule_TCL_Pbft() *TestModule_TCL_Pbft {
	return &TestModule_TCL_Pbft{
		epochID:           -1,
		totTxLatencyEpoch: make([]float64, 0),
		txNum:             make([]float64, 0),
	}
}

func (tml *TestModule_TCL_Pbft) OutputMetricName() string {
	return "Transaction_Confirm_Latency"
}

// modified latency
func (tml *TestModule_TCL_Pbft) UpdateMeasureRecord(b *pbft.BlockInfoMsg) {
	if b.BlockBodyLength == 0 { // empty block
		return
	}

	epochid := b.Epoch
	txs := b.ExcutedTxs
	mTime := b.CommitTime

	// extend
	for tml.epochID < epochid {
		tml.txNum = append(tml.txNum, 0)
		tml.totTxLatencyEpoch = append(tml.totTxLatencyEpoch, 0)
		tml.epochID++
	}

	for _, tx := range txs {
		if !tx.Time.IsZero() {
			tml.totTxLatencyEpoch[epochid] += mTime.Sub(tx.Time).Seconds()
			tml.txNum[epochid]++
		}
	}
}

func (tml *TestModule_TCL_Pbft) OutputRecord() (perEpochLatency []float64, totLatency float64) {
	perEpochLatency = make([]float64, 0)
	latencySum := 0.0
	totTxNum := 0.0
	for eid, totLatency := range tml.totTxLatencyEpoch {
		perEpochLatency = append(perEpochLatency, totLatency/tml.txNum[eid])
		latencySum += totLatency
		totTxNum += tml.txNum[eid]
	}
	totLatency = latencySum / totTxNum
	return
}

func (tml *TestModule_TCL_Pbft) GetDesc() metrics.Desc {
	b := metrics.NewDescBuilder("交易提交延迟(TCL)", "交易从到达到提交的延迟，即Tx Confirm Latency")

	var perEpochLatency []float64
	var totLatency float64

	perEpochLatency = make([]float64, 0)
	latencySum := 0.0
	totTxNum := 0.0
	for eid, totLatency := range tml.totTxLatencyEpoch {
		b.AddElem(fmt.Sprintf("第%v批次 交易提交延迟", eid+1), "各批次交易从到达到提交的延迟，即Tx Confirm Latency", fmt.Sprintf("%.2f s", totLatency/tml.txNum[eid]))
		perEpochLatency = append(perEpochLatency, totLatency/tml.txNum[eid])
		latencySum += totLatency
		totTxNum += tml.txNum[eid]
	}
	totLatency = latencySum / totTxNum
	b.AddElem("总计交易提交延迟", "整个过程中交易从到达到提交的延迟，即Tx Confirm Latency", fmt.Sprintf("%.2f s", totLatency))
	return b.GetDesc()
}
