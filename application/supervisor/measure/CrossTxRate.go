package measure

import (
	"fmt"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/metrics"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
)

// to test cross-transaction rate
type TestCrossTxRate_Relay struct {
	epochID       int
	totTxNum      []float64
	totCrossTxNum []float64
	relayTxRecord map[string]bool // record whether the relay1 has counted
}

func NewTestCrossTxRate_Relay() *TestCrossTxRate_Relay {
	return &TestCrossTxRate_Relay{
		epochID:       -1,
		totTxNum:      make([]float64, 0),
		totCrossTxNum: make([]float64, 0),
		relayTxRecord: make(map[string]bool),
	}
}

func (tctr *TestCrossTxRate_Relay) OutputMetricName() string {
	return "CrossTransaction_ratio"
}

func (tctr *TestCrossTxRate_Relay) UpdateMeasureRecord(b *pbft.BlockInfoMsg) {
	if b.BlockBodyLength == 0 { // empty block
		return
	}
	epochid := b.Epoch
	// extend
	for tctr.epochID < epochid {
		tctr.totTxNum = append(tctr.totTxNum, 0)
		tctr.totCrossTxNum = append(tctr.totCrossTxNum, 0)
		tctr.epochID++
	}

	// add relay1 txs
	// modify the relay map
	for _, r1tx := range b.Relay1Txs {
		tctr.relayTxRecord[string(r1tx.Hash)] = true
		tctr.totCrossTxNum[epochid] += 0.5
		tctr.totTxNum[epochid] += 0.5
	}
	// add inner-shard transaction and relay2 transactions
	for _, tx := range b.ExcutedTxs {
		if _, ok := tctr.relayTxRecord[string(tx.Hash)]; !ok {
			// inner-shard transaction
			tctr.totTxNum[epochid] += 1
		} else {
			tctr.totTxNum[epochid] += 0.5
			tctr.totCrossTxNum[epochid] += 0.5
		}
	}
}

func (tctr *TestCrossTxRate_Relay) OutputRecord() (perEpochCTXratio []float64, totCTXratio float64) {
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

func (tctr *TestCrossTxRate_Relay) GetDesc() metrics.Desc {
	b := metrics.NewDescBuilder("TPS avg", "平均每秒产生的交易，衡量交易的次数。单位为 交易/秒")

	var perEpochCTXratio []float64

	perEpochCTXratio = make([]float64, 0)
	allEpoch_totTxNum := 0.0
	allEpoch_ctxNum := 0.0
	for eid, totTxN := range tctr.totTxNum {
		b.AddElem(fmt.Sprintf("Epoch %v", eid), "", tctr.totCrossTxNum[eid]/totTxN)
		perEpochCTXratio = append(perEpochCTXratio, tctr.totCrossTxNum[eid]/totTxN)
		allEpoch_totTxNum += totTxN
		allEpoch_ctxNum += tctr.totCrossTxNum[eid]
	}
	perEpochCTXratio = append(perEpochCTXratio, allEpoch_totTxNum)
	perEpochCTXratio = append(perEpochCTXratio, allEpoch_ctxNum)
	b.AddElem("All Epoch Total Tx Num", "", allEpoch_totTxNum)
	b.AddElem("All Epoch Ctx Num", "", allEpoch_ctxNum)

	b.AddElem("Total CTX Ratio", "", allEpoch_ctxNum/allEpoch_totTxNum)
	return b.GetDesc()
}
