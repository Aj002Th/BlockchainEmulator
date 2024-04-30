package meter

import (
	"time"

	"github.com/Aj002Th/BlockchainEmulator/signal"
)

var Time_Begin time.Time
var time_result time.Duration

func onBegin(Nothing) {
	Time_Begin = time.Now()
}

func onStop(Nothing) {
	time_result = time.Since(Time_Begin)
}

func StartTimeCnt() {
	signal.GetSignalByName[Nothing]("OnNodeStart").Connect(onBegin)
	signal.GetSignalByName[Nothing]("OnNodeStop").Connect(onStop)
}
