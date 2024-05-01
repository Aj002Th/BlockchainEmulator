package network

// Network 实现p2p网络通信
type Network interface {
	Serve(endpoint string) chan []byte
	Send(context []byte, addr string)
	Broadcast(sender string, receivers []string, msg []byte)
	Receive()
	Close()
}
