package base

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"
)

// Account 预留结构, 便于以后扩展
// 平台的交互过程暂不涉及客户端加解密的操作
type Account struct {
	AccountAddress string
	PublicKey      []byte
}

// AccountState 账户状态详情
type AccountState struct {
	AccountAddress string
	Nonce          uint64
	Balance        *big.Int // 账户余额
}

// Deduct 扣款, 减少一个账户的余额
func (as *AccountState) Deduct(val *big.Int) bool {
	if as.Balance.Cmp(val) < 0 {
		return false
	}
	as.Balance.Sub(as.Balance, val)
	return true
}

// Deposit 存款, 增加一个账户的余额
func (as *AccountState) Deposit(value *big.Int) {
	as.Balance.Add(as.Balance, value)
}

// Hash 对AccountState求哈希
func (as *AccountState) Hash() []byte {
	h := sha256.Sum256(as.Encode())
	return h[:]
}

// Encode 对AccountState进行编码
func (as *AccountState) Encode() []byte {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(as)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// DecodeToAccountState 对AccountState进行解码
func DecodeToAccountState(b []byte) *AccountState {
	var as AccountState
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&as)
	if err != nil {
		log.Panic(err)
	}
	return &as
}
