package meter

func NodeSideStart() {
	// Node需要启动区块计数等等模块。
	// signal.GetSignalByName[Void]("OnNodeStart").Connect(func(Void) {
	StartNet()
	StartTimeCnt()
	StartPs()
	// })
}
