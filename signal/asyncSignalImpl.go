package signal

import (
	"log"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/supervisor_log"
)

var log1 *log.Logger = supervisor_log.Log1

// Channel和FANOUT风格的异步信号机制。
// 只能支持在进程内通信。

type ChanCancel = chan int

type Val[DATA any] struct {
	cd chan DATA
	cc ChanCancel
}

func NewVal[DATA any]() Val[DATA] {
	var cd chan DATA = make(chan DATA)
	var cc ChanCancel = make(ChanCancel)
	return Val[DATA]{cd: cd, cc: cc}
}

type AsyncSignalImpl[DATA any] struct {
	name        string
	outChannels map[*func(data DATA)]Val[DATA] // 从Handler到Val（cd和cc）
}

func NewAsyncSignalImpl[DATA any](name string) AsyncSignalImpl[DATA] {
	return AsyncSignalImpl[DATA]{
		name:        name,
		outChannels: make(map[*func(data DATA)]Val[DATA]),
	}
}

func (self AsyncSignalImpl[DATA]) Connect(cb func(data DATA)) bool { // 到时候只准传函数和lambda，不准传方法。

	val := NewVal[DATA]()
	self.outChannels[&cb] = val
	go func() { // 运行消息队列
		for {
			log1.Printf("I'm %v, now waiting for channel %v\n", self.name, val.cd)
			select {
			case <-val.cc:
				return
			case data := <-val.cd:
				cb(data)
			}
		}
	}()
	return true
}

// 返回值：操作是否成功
func (self AsyncSignalImpl[DATA]) Disconnect(cb func(data DATA)) bool {
	self.outChannels[&cb].cc <- 1 // 发送Cancel Token
	delete(self.outChannels, &cb)
	return true
}

func (self AsyncSignalImpl[DATA]) Emit(data DATA) {
	for _, val := range self.outChannels {
		log1.Printf("I'm %v, sending to channel %v\n", self.name, val.cd)
		val.cd <- data
	}
}
