package meter

func NodeSideStart() {
	// Node需要启动区块计数等等模块。
	StartNet()
	StartTimeCnt()
	StartPs()
}
