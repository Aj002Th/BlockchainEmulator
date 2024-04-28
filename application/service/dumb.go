package service

import "github.com/Aj002Th/BlockchainEmulator/application/comm"

type Dumb struct {
}

func (me *Dumb) Start() {
	// do nothing
}

func (me *Dumb) Gather(m *comm.MM) {
	// do nothing
}
