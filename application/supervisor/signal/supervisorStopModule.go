package signal

import (
	"sync"
)

// StopSignal to judge when the listener to send the stop message to the leaders
type StopSignal struct {
	stopLock      sync.Mutex // check the stopGap will not be modified by other processes``
	stopGap       int        // record how many empty txLists from leaders in a row
	stopThreshold int        // the threshold
}

func NewStopSignal(stop_Threshold int) *StopSignal {
	return &StopSignal{
		stopGap:       0,
		stopThreshold: stop_Threshold,
	}
}

// StopGapInc when receiving a message with an empty txList, then call this function to increase stopGap
func (ss *StopSignal) StopGapInc() {
	ss.stopLock.Lock()
	defer ss.stopLock.Unlock()
	ss.stopGap++
}

// StopGapReset when receiving a message with txs excuted, then call this function to reset stopGap
func (ss *StopSignal) StopGapReset() {
	ss.stopLock.Lock()
	defer ss.stopLock.Unlock()
	ss.stopGap = 0
}

// GapEnough Check the stopGap is enough or not
// if StopGap is not less than stopThreshold, then the stop message should be sent to leaders.
func (ss *StopSignal) GapEnough() bool {
	ss.stopLock.Lock()
	defer ss.stopLock.Unlock()
	return ss.stopGap >= ss.stopThreshold
}
