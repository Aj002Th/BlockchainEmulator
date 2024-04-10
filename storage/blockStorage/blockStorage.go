package blockStorage

import "github.com/Aj002Th/BlockchainEmulator/data/base"

// BlockStorage 针对区块的持久化存储
type BlockStorage interface {
	UpdateNewestBlock(newestBlockHash []byte)
	AddBlockHeader(blockHash []byte, bh *base.BlockHeader)
	AddBlock(block *base.Block)
	GetBlockHeader(blockHash []byte) (*base.BlockHeader, error)
	GetBlock(blockHash []byte) (*base.Block, error)
	GetNewestBlockHash() ([]byte, error)
	Close() error
}
