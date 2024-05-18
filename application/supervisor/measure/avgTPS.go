package measure

import (
	"fmt"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/metrics"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
)

type AvgTPS struct {
	epochID       int
	executedTxNum []float64
	startTime     []time.Time
	endTime       []time.Time
}

func NewAvgTPS() *AvgTPS {
	return &AvgTPS{
		epochID:       -1,
		executedTxNum: make([]float64, 0),
		startTime:     make([]time.Time, 0),
		endTime:       make([]time.Time, 0),
	}
}

func (tat *AvgTPS) OutputMetricName() string {
	return "Average_TPS"
}

func (tat *AvgTPS) UpdateMeasureRecord(b *pbft.BlockInfoMsg) {
	if b.BlockBodyLength == 0 {
		return
	}

	epochid := b.Epoch
	earliestTime := b.ProposeTime
	latestTime := b.CommitTime

	for tat.epochID < epochid {
		tat.executedTxNum = append(tat.executedTxNum, 0)
		tat.startTime = append(tat.startTime, time.Time{})
		tat.endTime = append(tat.endTime, time.Time{})
		tat.epochID++
	}
	for _, _ = range b.ExcutedTxs {
		tat.executedTxNum[epochid] += 1

	}
	if tat.startTime[epochid].IsZero() || tat.startTime[epochid].After(earliestTime) {
		tat.startTime[epochid] = earliestTime
	}
	if tat.endTime[epochid].IsZero() || latestTime.After(tat.endTime[epochid]) {
		tat.endTime[epochid] = latestTime
	}
}

func (tat *AvgTPS) OutputRecord() (perEpochTPS []float64, totalTPS float64) {
	perEpochTPS = make([]float64, tat.epochID+1)
	totalTxNum := 0.0
	eTime := time.Now()
	lTime := time.Time{}
	for eid, exTxNum := range tat.executedTxNum {
		timeGap := tat.endTime[eid].Sub(tat.startTime[eid]).Seconds()
		perEpochTPS[eid] = exTxNum / timeGap
		totalTxNum += exTxNum
		if eTime.After(tat.startTime[eid]) {
			eTime = tat.startTime[eid]
		}
		if tat.endTime[eid].After(lTime) {
			lTime = tat.endTime[eid]
		}
	}
	totalTPS = totalTxNum / (lTime.Sub(eTime).Seconds())
	return
}

func fmtTxPerSec(v float64) string {
	return fmt.Sprintf("%.2f txs/s", v)
}

func (tat *AvgTPS) GetDesc() metrics.Desc {

	b := metrics.NewDescBuilder("交易共识频率(AverageTPS)", "平均每秒产生的交易，衡量交易的次数。单位为 交易/秒")

	var perEpochTPS []float64
	var totalTPS float64
	perEpochTPS = make([]float64, tat.epochID+1)
	totalTxNum := 0.0
	eTime := time.Now()
	lTime := time.Time{}
	for eid, exTxNum := range tat.executedTxNum {
		timeGap := tat.endTime[eid].Sub(tat.startTime[eid]).Seconds()
		b.AddElem(fmt.Sprintf("第%v批次 交易共识频率", eid+1), fmt.Sprintf("第%v批次过程中产生交易的交易共识频率", eid+1), fmtTxPerSec(exTxNum/timeGap))
		perEpochTPS[eid] = exTxNum / timeGap
		totalTxNum += exTxNum
		if eTime.After(tat.startTime[eid]) {
			eTime = tat.startTime[eid]
		}
		if tat.endTime[eid].After(lTime) {
			lTime = tat.endTime[eid]
		}
	}
	totalTPS = totalTxNum / (lTime.Sub(eTime).Seconds())
	b.AddElem("总交易共识频率", "整个过程中平均每秒产生的交易", fmtTxPerSec(totalTPS))
	return b.GetDesc()

}
