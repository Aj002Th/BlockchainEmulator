package meter

import (
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/measure"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
	"github.com/Aj002Th/BlockchainEmulator/signal"
)

// 包括TPS，TCL那些

var m1, m2, m3 measure.MeasureModule

// 在sup端计算的时候用。
func CommitFeed(bim pbft.BlockInfoMsg) {
	m1.UpdateMeasureRecord(&bim)
	m2.UpdateMeasureRecord(&bim)
	m3.UpdateMeasureRecord(&bim)
}

func StartCommitRelate() {
	m1 = measure.NewTestModule_TCL_Relay()
	m2 = measure.NewTestModule_avgTPS_Relay()
	m3 = measure.NewPCL()
	sig := signal.GetSignalByName[pbft.BlockInfoMsg]("OnBimReached")
	sig.Connect(CommitFeed)
}
