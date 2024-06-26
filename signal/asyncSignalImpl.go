package signal

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/supervisor_log"
	"github.com/Aj002Th/BlockchainEmulator/logger"
	"github.com/chebyrash/promise"
)

// Channel和FANOUT风格的异步信号机制。
// 只能支持在进程内通信。
// Emit可能阻塞。如果回调阻塞的话。但能够保持data调用的顺序，也能防止竞争。
// 所以不要让回调阻塞。

type ChanCancel = chan int

type Val[DATA any] struct {
	cd       chan DATA
	cc       ChanCancel
	canceled *atomic.Bool
}

func NewVal[DATA any]() Val[DATA] {
	var canceled atomic.Bool
	cd := make(chan DATA)
	cc := make(ChanCancel)
	return Val[DATA]{cd: cd, cc: cc, canceled: &canceled}
}

type AsyncSignalImpl[DATA any] struct { // SO说by value在小对象的时候开销更小。然后它是深拷贝，但遇到指针，map，chan之类就会停止往下拷贝。所以这样应该是可以的。
	name        string
	outChannels map[string]Val[DATA] // 从Handler到Val（cd和cc）
	lastEmit    *promise.Promise[int]
}

func NewAsyncSignalImpl[DATA any](name string) AsyncSignalImpl[DATA] {
	if ExistSignalName(name) {
		return GetSignalByName[DATA](name).(AsyncSignalImpl[DATA]) // 事实是只有这个Async的Signal有用。所以默认全是AsyncSignal。
	}
	s := AsyncSignalImpl[DATA]{
		name:        name,
		outChannels: make(map[string]Val[DATA]),
		lastEmit:    promise.New(func(resolve func(int), reject func(error)) { resolve(0) }),
	}
	RegisterSig[DATA](s) // 把它注册到全局的表里。方便稍后用名字查到。
	return s
}

func (self AsyncSignalImpl[DATA]) GetName() string {
	return self.name
}

// Connect cb在一个goroutine上依次调用，没有数据竞争。所以cb不要写阻塞代码。
func (self AsyncSignalImpl[DATA]) Connect(cb func(data DATA)) bool { // 到时候只准传函数和lambda，不准传方法。
	logger.Printf("Connect &cb: %v\n", &cb)
	logger.Printf("Connect cb: %v\n", cb)
	val := NewVal[DATA]()
	logger.Printf("val: %v\n", val)
	self.outChannels[CbToIdx(cb)] = val
	go func() { // 运行消息队列
		for {
			logger.Printf("Signal %v, now waiting for channel %v\n", self.name, val.cd)

			data := <-val.cd
			if val.canceled.Load() {
				return
			}
			logger.Printf("Signal %v, now calling callback from channel %v\n", self.name, val.cd)
			cb(data)
		}
	}()
	return true
}

// Disconnect 返回值：操作是否成功
func (self AsyncSignalImpl[DATA]) Disconnect(cb func(data DATA)) bool {
	supervisor_log.DebugLog.Printf("Signal %v, now setting cancel\n", self.name)
	supervisor_log.DebugLog.Printf("self.outChannels: %v\n", self.outChannels)
	supervisor_log.DebugLog.Printf("Disc &cb: %v\n", &cb)
	supervisor_log.DebugLog.Printf("Disc cb: %v\n", cb)
	self.outChannels[CbToIdx(cb)].canceled.Store(true)
	delete(self.outChannels, CbToIdx(cb))
	return true
}

// Emit 不会阻塞, 用Promise Then串起来了, 同时能够保持data调用的顺序。
func (self AsyncSignalImpl[DATA]) Emit(data DATA) {
	for _, val := range self.outChannels {
		supervisor_log.DebugLog.Printf("Signal %v, sending to channel %v\n", self.name, val.cd)
		ctx := context.Background()
		self.lastEmit = promise.Then(self.lastEmit, ctx, func(int) (int, error) {
			val.cd <- data
			return 0, nil
		})

	}
}

func CbToIdx[DATA any](cb func(data DATA)) string {
	return fmt.Sprintf("%v", cb)
}
