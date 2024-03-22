package base

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

// BlockHeader 区块头
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

// Block 区块
type Block struct {
	Header *BlockHeader
	Body   []*Transaction
	Hash   []byte
}

// NewBlock 创建区块
// 调用前需要保证区块头中的信息完善且正确
// 创建后不应进行修改操作, 否则哈希结果出错
func NewBlock(header *BlockHeader, body []*Transaction) *Block {
	block := &Block{
		Header: header,
		Body:   body,
	}
	block.Hash = block.Header.Hash()
	return block
}

// Encode 对Block进行编码
func (b *Block) Encode() []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// DecodeBlock 对Block进行解码
func DecodeBlock(b []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}
