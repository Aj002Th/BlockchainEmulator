package meter

import (
	"os"
	"sync/atomic"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/signal"
	"github.com/shirou/gopsutil/process"
)

var AvgCpuPercent float64

var cpuSampleCnt int
var DiskMetric uint64
var diskSampleCnt int

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
		for {
			time.Sleep(1000)
			if stop.Load() {
				return
			}
			i, err := p.CPUPercent()
			t := i
			if err != nil {
				panic("Wrong")
			}
			// AvgCpuTime = t
			AvgCpuPercent = (AvgCpuPercent*float64(cpuSampleCnt) + t) / (float64(cpuSampleCnt) + 1)
			cpuSampleCnt++
		}
	}()

	go func() {
		for {
			time.Sleep(1000)
			if stop.Load() {
				return
			}
			i, err := p.MemoryInfo()
			t := i.RSS
			if err != nil {
				panic("Wrong")
			}
			DiskMetric = (DiskMetric*uint64(diskSampleCnt) + t) / uint64(diskSampleCnt+1)
			diskSampleCnt++
		}
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
