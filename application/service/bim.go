package service

import (
	"github.com/Aj002Th/BlockchainEmulator/application/comm"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
	"github.com/Aj002Th/BlockchainEmulator/signal"
)

type Bim struct {
	c (chan pbft.BlockInfoMsg)
}

func NewBim() Bim {
	c := make(chan pbft.BlockInfoMsg)
	return Bim{c: c}
}

func (me *Bim) Start() {
	sig := signal.FindSignalByName[pbft.BlockInfoMsg]("OnBlockInfoMsg")
	sig.Connect(func(data pbft.BlockInfoMsg) {
		go func() { me.c <- data }() // 别把signal阻塞了。
	})

}

func (me *Bim) Gather(m *comm.MM) { // 用生产者消费者把它稍微适配一下。
	m.Bim = <-me.c
}
