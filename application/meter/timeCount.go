package meter

import "time"

var time_begin time.Time
var time_result time.Duration

func onBegin() {
	time_begin = time.Now()
}

func onEnd() {
	time_result = time.Since(time_begin)
}
