package chain

import (
	"github.com/Aj002Th/BlockchainEmulator/storage/blockStorage"
	"github.com/Aj002Th/BlockchainEmulator/storage/stateStorage"
)

type Option func(bc *BlockChain)

func WithStateStorageLevelDB(path string) Option {
	return func(bc *BlockChain) {
		bc.StateStorage, _ = stateStorage.NewLevelDB(path)
	}
}

func WithStateStorageMemKVStore() Option {
	return func(bc *BlockChain) {
		bc.StateStorage = stateStorage.NewMemKVStore()
	}
}

func WithBlockStorageBolt(nodeID uint) Option {
	return func(bc *BlockChain) {
		bc.BlockStorage = blockStorage.NewBoltStorage(nodeID)
	}
}
