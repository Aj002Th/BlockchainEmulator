package blockStorage

import (
	"github.com/Aj002Th/BlockchainEmulator/data/base"
)

type RedisStorage struct {
	dbFilePath            string // redis 数据库文件路径
	blockBucket           string // redis 中存放 block 数据的 key prefix
	blockHeaderBucket     string // redis 中存放 block header 数据的 key prefix
	newestBlockHashBucket string // redis 中存放 block 数据的 key prefix
}

func (s *RedisStorage) UpdateNewestBlock(newestBlockHash []byte) {
	//TODO implement me
	panic("implement me")
}

func (s *RedisStorage) AddBlockHeader(blockHash []byte, bh *base.BlockHeader) {
	//TODO implement me
	panic("implement me")
}

func (s *RedisStorage) AddBlock(block *base.Block) {
	//TODO implement me
	panic("implement me")
}

func (s *RedisStorage) GetBlockHeader(blockHash []byte) (*base.BlockHeader, error) {
	//TODO implement me
	panic("implement me")
}

func (s *RedisStorage) GetBlock(blockHash []byte) (*base.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (s *RedisStorage) GetNewestBlockHash() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func NewRedisStorage(nodeID uint) *RedisStorage {
	return &RedisStorage{}
}
