package blockStorage

import "github.com/Aj002Th/BlockchainEmulator/data/base"

type RedisStorage struct {
}

func (r RedisStorage) UpdateNewestBlock(newestBlockHash []byte) {
	//TODO implement me
	panic("implement me")
}

func (r RedisStorage) AddBlockHeader(blockHash []byte, bh *base.BlockHeader) {
	//TODO implement me
	panic("implement me")
}

func (r RedisStorage) AddBlock(block *base.Block) {
	//TODO implement me
	panic("implement me")
}

func (r RedisStorage) GetBlockHeader(blockHash []byte) (*base.BlockHeader, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisStorage) GetBlock(blockHash []byte) (*base.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisStorage) GetNewestBlockHash() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (r RedisStorage) Close() error {
	//TODO implement me
	panic("implement me")
}
