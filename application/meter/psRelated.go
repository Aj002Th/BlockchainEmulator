package meter

import (
	"os"
	"time"

	"github.com/shirou/gopsutil/process"
)

type Ctx struct {
	avgCpuTime float64
	diskMetric uint64
}

func Start(ctx Ctx) {
	// 创建统计进程
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		panic("")
	}

	go func() {
		time.Sleep(1000)
		i, err := p.Times()
		t := i.Total()
		if err != nil {
			panic("Wrong")
		}
		ctx.avgCpuTime = t
	}()

	go func() {
		time.Sleep(1000)
		i, err := p.MemoryInfo()
		t := i.RSS
		if err != nil {
			panic("Wrong")
		}
		ctx.diskMetric = t
	}()

}
