package signal

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/supervisor_log"
	"github.com/chebyrash/promise"
)

var log1 *log.Logger = supervisor_log.Log1

// Channel和FANOUT风格的异步信号机制。
// 只能支持在进程内通信。
// Emit可能阻塞。如果回调阻塞的话。但能够保持data调用的顺序，也能防止竞争。
// 所以不要让回调阻塞。

type ChanCancel = chan int

type Val[DATA any] struct { // TODO: 加一个调试信息，比如函数名字之类，方便调试。
	cd       chan DATA
	cc       ChanCancel
	canceled *atomic.Bool
}

func NewVal[DATA any]() Val[DATA] {
	var cd chan DATA = make(chan DATA)
	var cc ChanCancel = make(ChanCancel)
	var canceled atomic.Bool
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
	RegisterSig(s) // 把它注册到全局的表里。方便稍后用名字查到。
	return s
}

func (self AsyncSignalImpl[DATA]) GetName() string {
	return self.name
}

// cb在一个goroutine上依次调用，没有数据竞争。所以cb不要写阻塞代码。
func (self AsyncSignalImpl[DATA]) Connect(cb func(data DATA)) bool { // 到时候只准传函数和lambda，不准传方法。
	fmt.Printf("Connect &cb: %v\n", &cb)
	fmt.Printf("Connect cb: %v\n", cb)
	val := NewVal[DATA]()
	fmt.Printf("val: %v\n", val)
	self.outChannels[CbToIdx(cb)] = val
	go func() { // 运行消息队列
		for {
			log1.Printf("Signal %v, now waiting for channel %v\n", self.name, val.cd)

			data := <-val.cd
			if val.canceled.Load() {
				return
			}
			log1.Printf("Signal %v, now calling callback from channel %v\n", self.name, val.cd)
			cb(data)
		}
	}()
	return true
}

// 返回值：操作是否成功
func (self AsyncSignalImpl[DATA]) Disconnect(cb func(data DATA)) bool {
	log1.Printf("Signal %v, now setting cancel\n", self.name)
	fmt.Printf("self.outChannels: %v\n", self.outChannels)
	fmt.Printf("Disc &cb: %v\n", &cb)
	fmt.Printf("Disc cb: %v\n", cb)
	self.outChannels[CbToIdx(cb)].canceled.Store(true)
	delete(self.outChannels, CbToIdx(cb))
	return true
}

// 不会阻塞，用Promise Then串起来了。但能够保持data调用的顺序。
func (self AsyncSignalImpl[DATA]) Emit(data DATA) {
	for _, val := range self.outChannels {
		log1.Printf("Signal %v, sending to channel %v\n", self.name, val.cd)
		// val.cd <- data
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
