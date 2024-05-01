package network

import "github.com/Aj002Th/BlockchainEmulator/signal"

// Network 实现p2p网络通信
type Network interface {
	Serve(endpoint string) chan []byte
	Send(context []byte, addr string)
	Broadcast(sender string, receivers []string, msg []byte)
	GetOnUpload() signal.Signal[int]
	GetOnDownload() signal.Signal[int]
	Close()
}

var Tcp Network = NewTcpCustomProtocolNetwork()
