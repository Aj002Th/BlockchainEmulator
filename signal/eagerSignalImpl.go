// 从https://github.com/vorot93/golang-signals/blob/master/signals.go直接摘来的，做了替换，直接以cb为索引。

package signal

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
