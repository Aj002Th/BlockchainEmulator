package meter

import (
	"time"

	"github.com/Aj002Th/BlockchainEmulator/signal"
)

var time_begin time.Time
var time_result time.Duration

func onBegin(Nothing) {
	time_begin = time.Now()
}

func onStop(Nothing) {
	time_result = time.Since(time_begin)
}

func StartTimeCnt() {
	signal.GetSignalByName[Nothing]("EmulatorOnBegin").Connect(onBegin)
	signal.GetSignalByName[Nothing]("EmulatorOnStop").Connect(onStop)
}
