package async_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/signal"
)

func toEmit1() chan int {
	var completion chan int = make(chan int)
	var sig signal.Signal[int] = signal.NewAsyncSignalImpl[int]("sigTest")
	sig.Connect(func(data int) {
		fmt.Println("there is cb")
		println("the data is : ", data)
		completion <- 1
	})
	sig.Emit(1)
	return completion
}

func toEmitMulti(n int) chan int {
	var completion chan int = make(chan int)
	var sig signal.Signal[int] = signal.NewAsyncSignalImpl[int]("sigTest")
	for i := 0; i < n; i++ {

		sig.Connect(func(data int) {
			fmt.Println("there is cb ", i)
			println("the data is : ", data)
			completion <- 1
		})
	}
	sig.Emit(1)
	return completion
}

func Test1(t *testing.T) {
	var completion = toEmit1()
	select {
	case <-completion:
		return
	case <-time.After(2 * time.Second):
		t.Fatalf(`哦不， signal handler 没有反应！`)
	}
}

func TestMulti(t *testing.T) {
	var completion = toEmitMulti(10)
	select {
	case <-completion:
		return
	case <-time.After(2 * time.Second):
		t.Fatalf(`哦不， signal handler 没有反应！`)
	}
}

func TestDisconnect(t *testing.T) {
	var sig signal.Signal[int] = signal.NewAsyncSignalImpl[int]("sigTest")
	cnt := 0
	var cb func(data int)
	cb = func(data int) {
		fmt.Println("there is cb")
		println("the data is : ", data)
		cnt++
		sig.Disconnect(cb)
	}
	sig.Connect(cb)
	sig.Emit(1)
	sig.Emit(2)
	time.Sleep(1 * time.Second)
	if cnt != 1 {
		t.Fatalf("")
	}
}
