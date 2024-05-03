package network

import (
	"github.com/Aj002Th/BlockchainEmulator/signal"
)

type HttpNetwork struct {
}

func (h *HttpNetwork) Serve(endpoint string) chan []byte {
	//TODO implement me
	panic("implement me")
}

func (h *HttpNetwork) Send(context []byte, addr string) {
	//TODO implement me
	panic("implement me")
}

func (h *HttpNetwork) Broadcast(sender string, receivers []string, msg []byte) {
	//TODO implement me
	panic("implement me")
}

func (h *HttpNetwork) GetOnUpload() signal.Signal[int] {
	//TODO implement me
	panic("implement me")
}

func (h *HttpNetwork) GetOnDownload() signal.Signal[int] {
	//TODO implement me
	panic("implement me")
}

func (h *HttpNetwork) Close() {
	//TODO implement me
	panic("implement me")
}
