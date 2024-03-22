package data

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type BlockHeader struct {
	ParentBlockHash []byte    // 父区块哈希
	StateRoot       []byte    // 状态树根哈希
	TxRoot          []byte    // 交易的哈希
	Number          uint64    // 区块中包含的交易数量
	Time            time.Time // 区块创建时间
	Miner           uint64    // 区块创建节点
}

// Encode 对BlockHeader进行编码
func (bh *BlockHeader) Encode() []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(bh)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// DecodeBlockHeader 对BlockHeader进行解码
func DecodeBlockHeader(b []byte) *BlockHeader {
	var blockHeader BlockHeader
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&blockHeader)
	if err != nil {
		log.Panic(err)
	}
	return &blockHeader
}

// Hash 对BlockHeader求哈希
func (bh *BlockHeader) Hash() []byte {
	hash := sha256.Sum256(bh.Encode())
	return hash[:]
}

type Block struct {
}
