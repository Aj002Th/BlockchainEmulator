package base

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"
	"time"
)

// Transaction 交易
type Transaction struct {
	Sender    string
	Recipient string
	Nonce     uint64
	Value     *big.Int
	Hash      []byte

	// 一些额外信息
	Time time.Time // 添加到交易池的时间
}

// NewTransaction 创建交易
// 创建后不应进行修改操作, 否则哈希结果出错
func NewTransaction(sender, recipient string, value *big.Int, nonce uint64) *Transaction {
	tx := &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Value:     value,
		Nonce:     nonce,
	}
	hash := sha256.Sum256(tx.Encode())
	tx.Hash = hash[:]
	return tx
}

// Encode 对Transaction进行编码
func (tx *Transaction) Encode() []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// DecodeTransaction 对Transaction进行解码
func DecodeTransaction(b []byte) *Transaction {
	var tx Transaction
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}
	return &tx
}
