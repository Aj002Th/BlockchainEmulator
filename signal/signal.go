package signal

// 信号。实例代表某种信号。可以注册或者卸载lambda函数。
type Signal[DATA any] interface {
	Connect(cb func(data DATA)) bool
	Disconnect(cb func(data DATA)) bool
	Emit(data DATA)
}
