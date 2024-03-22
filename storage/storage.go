package storage

import "github.com/Aj002Th/BlockchainEmulator/data/base"

// Storage 储存层对外提供的能力
type Storage interface {
	UpdateNewestBlock(newestBlockHash []byte)
	AddBlockHeader(blockHash []byte, bh *base.BlockHeader)
	AddBlock(block *base.Block)
	GetBlockHeader(blockHash []byte) (*base.BlockHeader, error)
	GetBlock(blockHash []byte) (*base.Block, error)
	GetNewestBlockHash() ([]byte, error)
}
