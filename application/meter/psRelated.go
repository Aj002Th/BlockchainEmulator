package meter

import (
	"os"
	"sync/atomic"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/signal"
	"github.com/shirou/gopsutil/process"
)

var avgCpuTime float64
var diskMetric uint64

type Nothing = struct{}
type Void = struct{}

// 这个不用依赖信号，反正自力更生。
func StartPs() {
	// 创建统计进程
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		panic("")
	}
	var stop atomic.Bool
	sig := signal.GetSignalByName[Nothing]("OnEmulatorStop")
	sig.Connect(func(data Nothing) { stop.Store(true) })

	go func() {
		time.Sleep(1000)
		if stop.Load() {
			return
		}
		i, err := p.Times()
		t := i.Total()
		if err != nil {
			panic("Wrong")
		}
		avgCpuTime = t
	}()

	go func() {
		time.Sleep(1000)
		if stop.Load() {
			return
		}
		i, err := p.MemoryInfo()
		t := i.RSS
		if err != nil {
			panic("Wrong")
		}
		diskMetric = t
	}()

}

func check(stop chan Nothing) bool {
	select {
	case <-stop:
		return true
	default:
		return false
	}
}
