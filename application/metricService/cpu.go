package metricService

import (
	"os"
	"time"

	"github.com/shirou/gopsutil/process"
)

type Cpu struct {
	avg float64
	p   *process.Process
}

func (me *Cpu) Start() {
	// 创建统计进程
	var err error
	me.p, err = process.NewProcess(int32(os.Getpid()))
	if err != nil {
		panic("")
	}

	go func() {
		time.Sleep(1000)
		i, err := me.p.Times()
		t := i.Total()
		if err != nil {
			panic("Wrong")
		}
		me.avg = t
	}()
}

func (me *Cpu) GetCpuStatic() float64 {
	return me.avg
}

func (me *Cpu) Gather(m *MyInfo) {

}

// gopsutil是 Python 工具库psutil 的 Golang 移植版，可以帮助我们方便地获取各种系统和硬件信息。
// gopsutil为我们屏蔽了各个系统之间的差异，具有非常强悍的可移植性。
// 有了gopsutil，我们不再需要针对不同的系统使用syscall调用对应的系统方法。
// 更棒的是gopsutil的实现中没有任何cgo的代码，使得交叉编译成为可能。

// ---------------------------------------------------
