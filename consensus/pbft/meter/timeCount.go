package meter

import (
	"time"
)

var TimeBegin time.Time

func onBegin() {
	TimeBegin = time.Now()
}

func StartTimeCnt() {
	onBegin()
}
