package metricService

import (
	"encoding/json"

	"github.com/Aj002Th/BlockchainEmulator/params"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/push"
)

var s mangos.Socket

func startService() {
	var err error
	s, err = push.NewSocket()
	if err != nil {
		panic("new sock failed")
	}
	err = s.Dial(params.SupervisorEndpoint)
	if err != nil {
		panic("dial failed")
	}
}

func sendMsg(m MyInfo) {
	bs, err := json.Marshal(m)
	if err != nil {
		panic("json failed")
	}
	err = s.Send(bs)
	if err != nil {
		panic("send failed")
	}
}
