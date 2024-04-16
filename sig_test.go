package main

import (
	"testing"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/signal"
)

// Channel和FANOUT风格的异步信号机制

type ChanCancel = chan int

type Val[DATA any] struct {
	cd *chan DATA
	cc *ChanCancel
}

func NewVal[DATA any]() Val[DATA] {
	var cd chan DATA
	var cc ChanCancel
	return Val[DATA]{cd: &cd, cc: &cc}
}

type AsyncSignalImpl[DATA any] struct {
	name        string
	outChannels map[*func(data DATA)]*Val[DATA] // 从Handler到Val（cd和cc）

	www chan int
}

func NewAsyncSignalImpl[DATA any](name string) AsyncSignalImpl[DATA] {
	return AsyncSignalImpl[DATA]{
		name:        name,
		outChannels: make(map[*func(data DATA)]*Val[DATA]),
	}
}

func (self AsyncSignalImpl[DATA]) Connect(cb func(data DATA)) bool { // 到时候只准传函数和lambda，不准传方法。
	go func() {
		<-self.www
		var dd DATA
		cb(dd)
	}()
	// val := NewVal[DATA]()
	// self.outChannels[&cb] = &val
	// go func() { // 运行消息队列
	// 	for {
	// 		select {
	// 		case <-*val.cc:
	// 			return
	// 		case data := <-*val.cd:
	// 			cb(data)
	// 		}
	// 	}
	// }()
	return true
}

// 返回值：操作是否成功
func (self AsyncSignalImpl[DATA]) Disconnect(cb func(data DATA)) bool {
	*self.outChannels[&cb].cc <- 1 // 发送Cancel Token
	delete(self.outChannels, &cb)
	return true
}

func (self AsyncSignalImpl[DATA]) Emit(data DATA) {
	// for _, val := range self.outChannels {
	// 	*val.cd <- data
	// }
	self.www <- 1
}

func toEmit() chan int {
	var completion chan int
	var sig signal.Signal[int] = NewAsyncSignalImpl[int]("sigTest")
	sig.Connect(func(data int) {
		println(data)
		completion <- 1
	})
	sig.Emit(1)
	return completion
}

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {
	var completion = toEmit()
	select {
	case <-completion:
		return
	case <-time.After(2 * time.Second):
		t.Fatalf(`哦不， signal handler 没有反应！`)
	}
}

func Test1(t *testing.T) {
	var c chan int = make(chan int)
	go func() {
		<-c
	}()
	c <- 1
}
