package meter

import (
	"time"
)

var Time_Begin time.Time

func onBegin() {
	Time_Begin = time.Now()
}

func StartTimeCnt() {
	onBegin()
}
