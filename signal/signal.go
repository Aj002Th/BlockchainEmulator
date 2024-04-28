package signal

// 信号。实例代表某种信号。可以注册或者卸载lambda函数。
type Signal[DATA any] interface {
	GetName() string
	Connect(cb func(data DATA)) bool
	Disconnect(cb func(data DATA)) bool
	Emit(data DATA)
}

var globalSigs map[string]interface{}

func RegisterSig[DATA any](sig Signal[DATA]) {
	name := (sig).GetName() // 什么逆天，为什么要加星花？
	globalSigs[name] = sig
}

func FindSignalByName[DATA any](name string) Signal[DATA] {
	v, ok := globalSigs[name]
	if !ok {
		panic("error")
	}
	return v.(Signal[DATA])
}
