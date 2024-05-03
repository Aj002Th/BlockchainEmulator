package meter

import (
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/measure"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
	"github.com/Aj002Th/BlockchainEmulator/signal"
)

// 包括TPS，TCL那些
var mTCL, mAvgTPS, mPCL, mBlockNumCount, mTxNumCount measure.MeasureModule

// CommitFeed 在sup端计算的时候用。
func CommitFeed(bim *pbft.BlockInfoMsg) {
	mTCL.UpdateMeasureRecord(bim)
	mAvgTPS.UpdateMeasureRecord(bim)
	mPCL.UpdateMeasureRecord(bim)
	mBlockNumCount.UpdateMeasureRecord(bim)
	mTxNumCount.UpdateMeasureRecord(bim)
}

func StartCommitRelate() []*measure.MeasureModule {
	mTCL = measure.NewTCL()
	mAvgTPS = measure.NewAvgTPS()
	mPCL = measure.NewPCL()
	mBlockNumCount = measure.NewBlockNumCount()
	mTxNumCount = measure.NewTxNumCount()

	sig := signal.GetSignalByName[*pbft.BlockInfoMsg]("OnBimReached")
	sig.Connect(CommitFeed)
	return []*measure.MeasureModule{&mTCL, &mAvgTPS, &mPCL, &mBlockNumCount, &mTxNumCount}
}
