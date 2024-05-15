package signal

import (
	"sync"
)

type StopSignal struct {
	stopLock      sync.Mutex
	stopGap       int
	stopThreshold int
}

func NewStopSignal(stop_Threshold int) *StopSignal {
	return &StopSignal{
		stopGap:       0,
		stopThreshold: stop_Threshold,
	}
}

func (ss *StopSignal) StopGapInc() {
	ss.stopLock.Lock()
	defer ss.stopLock.Unlock()
	ss.stopGap++
}

func (ss *StopSignal) StopGapReset() {
	ss.stopLock.Lock()
	defer ss.stopLock.Unlock()
	ss.stopGap = 0
}

func (ss *StopSignal) GapEnough() bool {
	ss.stopLock.Lock()
	defer ss.stopLock.Unlock()
	return ss.stopGap >= ss.stopThreshold
}
