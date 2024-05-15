package signal

import "github.com/Aj002Th/BlockchainEmulator/logger"

// 信号。实例代表某种信号。可以注册或者卸载lambda函数。
type Signal[DATA any] interface {
	GetName() string
	Connect(cb func(data DATA)) bool
	Disconnect(cb func(data DATA)) bool
	Emit(data DATA)
}

// 不要直接操作这个。因为单例模式可能没初始化。请用GetGSig访问。
var globalSigs map[string]interface{} = make(map[string]interface{})

// 线程不安全。访问请加锁。
func GetGSig() map[string]interface{} {
	if globalSigs == nil {
		globalSigs = make(map[string]interface{})
	}
	return globalSigs
}

func RegisterSig[DATA any](sig Signal[DATA]) {
	name := (sig).GetName()
	GetGSig()[name] = sig
}

func GetSignalByName[DATA any](name string) Signal[DATA] {
	v, ok := globalSigs[name]
	if !ok {
		newVal := NewAsyncSignalImpl[DATA](name)
		logger.Printf("Warning: GetSig Creating a New Sig In GetByName %v Because No Name Found", name)
		globalSigs[name] = newVal
		return newVal
	}
	return v.(Signal[DATA])
}

func ExistSignalName(name string) bool {
	_, ok := globalSigs[name]
	if !ok {
		return false
	}
	return true
}
