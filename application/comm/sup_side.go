package comm

import (
	"encoding/json"

	"github.com/Aj002Th/BlockchainEmulator/params"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pull"
)

var s1 mangos.Socket

func Listen() {
	var err error
	s1, err = pull.NewSocket()
	if err != nil {
		panic("new sock failed")
	}
	err = s1.Listen(params.SupervisorEndpoint)
	if err != nil {
		panic("listen failed")
	}
}

func Recv() *MM {
	bs, err := s1.Recv()
	if err != nil {
		panic("recv failed")
	}
	var m MM
	err = json.Unmarshal(bs, &m)
	return &m
}
