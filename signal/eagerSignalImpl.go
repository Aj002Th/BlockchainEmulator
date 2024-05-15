package signal

// 这一代码受到了：
// https://github.com/vorot93/golang-signals/blob/master/signals.go
// 启发，并且做了改进，直接以cb为索引。
// 以此实现了一个传统意义上的Signal。

import (
	"sync"
)

type EagerSignalImpl[DATA any] struct {
	m     sync.Mutex
	slots map[*func(data DATA)]int
}

func (self *EagerSignalImpl[DATA]) Connect(cb func(data DATA)) bool {
	self.m.Lock()
	defer self.m.Unlock()
	if self.slots == nil {
		self.slots = *new(map[*func(data DATA)]int)
	}
	self.slots[&cb] = 1
	return true
}

// 返回值：操作是否成功
func (self *EagerSignalImpl[DATA]) Disconnect(cb *func(data DATA)) bool {
	self.m.Lock()
	defer self.m.Unlock()
	if self.slots == nil {
		self.slots = *new(map[*func(data DATA)]int)
	}
	var _, exists = self.slots[cb]
	if !exists {
		return false
	}
	delete(self.slots, cb)
	return true
}

func (self *EagerSignalImpl[DATA]) Emit(data DATA) {
	for cb, _ := range self.slots {
		(*cb)(data)
	}
}
