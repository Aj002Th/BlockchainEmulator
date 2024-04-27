package measure

import "github.com/Aj002Th/BlockchainEmulator/consensus/pbft"

// to test cross-transaction rate
type TestCpuUsage struct {
	epochID int
	txNum   []float64
}

func NewCpuUsage() *TestCpuUsage {
	return &TestCpuUsage{
		epochID: -1,
		txNum:   make([]float64, 0),
	}
}

func (ttnc *TestCpuUsage) OutputMetricName() string {
	return "CPU_Usage"
}

func (ttnc *TestCpuUsage) UpdateMeasureRecord(b *pbft.BlockInfoMsg) {
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

func (ttnc *TestCpuUsage) OutputRecord() (perEpochCTXs []float64, totTxNum float64) {
	perEpochCTXs = make([]float64, 0)
	totTxNum = 0.0
	for _, tn := range ttnc.txNum {
		perEpochCTXs = append(perEpochCTXs, tn)
		totTxNum += tn
	}
	return perEpochCTXs, totTxNum
}

func (ttnc *TestCpuUsage) GetDesc() string {
	return "CPU占用率"
}
