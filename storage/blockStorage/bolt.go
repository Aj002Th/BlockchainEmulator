package blockStorage

import (
	"errors"
	"fmt"
	"github.com/Aj002Th/BlockchainEmulator/data/base"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

var (
	blockBucket           = "block"
	blockHeaderBucket     = "blockHeader"
	newestBlockHashBucket = "newestBlockHash"
	newestBlockHashKey    = "OnlyNewestBlock"
)

// BoltStorage 基于 bolt 实现的存储
type BoltStorage struct {
	dbFilePath            string // bolt数据库文件路径
	blockBucket           string // bolt中存放 block 数据的 bucket
	blockHeaderBucket     string // bolt中存放 block header 数据的 bucket
	newestBlockHashBucket string // bolt中存放 block 数据的 bucket
	DB                    *bolt.DB
}

// NewBoltStorage 创建BoltStorage
// nodeID指代节点的唯一ID
func NewBoltStorage(nodeID uint) *BoltStorage {
	// 保证存放数据库文件的路径存在
	_, errStat := os.Stat("./record")
	if os.IsNotExist(errStat) {
		errMkdir := os.Mkdir("./record", os.ModePerm)
		if errMkdir != nil {
			log.Panic(errMkdir)
		}
	} else if errStat != nil {
		log.Panic(errStat)
	}

	dbFilePath := fmt.Sprintf("./record/node_%d.db", nodeID)
	s := &BoltStorage{
		dbFilePath:            dbFilePath,
		blockBucket:           blockBucket,
		blockHeaderBucket:     blockHeaderBucket,
		newestBlockHashBucket: newestBlockHashBucket,
		DB:                    nil,
	}

	//初始化数据库
	db, err := bolt.Open(dbFilePath, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	// 创建 buckets
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(s.blockBucket))
		if err != nil {
			log.Panic("create blocksBucket failed")
		}

		_, err = tx.CreateBucketIfNotExists([]byte(s.blockHeaderBucket))
		if err != nil {
			log.Panic("create blockHeaderBucket failed")
		}

		_, err = tx.CreateBucketIfNotExists([]byte(s.newestBlockHashBucket))
		if err != nil {
			log.Panic("create newestBlockHashBucket failed")
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	s.DB = db
	return s
}

// UpdateNewestBlock 更新newestBlockHashBucket中存储的"最新区块hash"
func (s *BoltStorage) UpdateNewestBlock(newestBlockHash []byte) {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		nbhBucket := tx.Bucket([]byte(s.newestBlockHashBucket))
		// 该 bucket 中仅存储了这一个键值对
		err := nbhBucket.Put([]byte(newestBlockHashKey), newestBlockHash)
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// AddBlockHeader 存储区块头
func (s *BoltStorage) AddBlockHeader(blockHash []byte, bh *base.BlockHeader) {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		bhBucket := tx.Bucket([]byte(s.blockHeaderBucket))
		err := bhBucket.Put(blockHash, bh.Encode())
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// AddBlock 存储区块
func (s *BoltStorage) AddBlock(block *base.Block) {
	err := s.DB.Update(func(tx *bolt.Tx) error {
		bBucket := tx.Bucket([]byte(s.blockBucket))
		err := bBucket.Put(block.Hash, block.Encode())
		if err != nil {
			log.Panic(err)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	s.AddBlockHeader(block.Hash, block.Header)
	s.UpdateNewestBlock(block.Hash)
}

// GetBlockHeader 获取区块头
func (s *BoltStorage) GetBlockHeader(blockHash []byte) (*base.BlockHeader, error) {
	var result *base.BlockHeader
	err := s.DB.View(func(tx *bolt.Tx) error {
		bhBucket := tx.Bucket([]byte(s.blockHeaderBucket))
		bhEncoded := bhBucket.Get(blockHash)
		if bhEncoded == nil {
			return errors.New("the block is not existed")
		}
		result = base.DecodeBlockHeader(bhEncoded)
		return nil
	})
	return result, err
}

// GetBlock 获取区块
func (s *BoltStorage) GetBlock(blockHash []byte) (*base.Block, error) {
	var result *base.Block
	err := s.DB.View(func(tx *bolt.Tx) error {
		bBucket := tx.Bucket([]byte(s.blockBucket))
		bEncoded := bBucket.Get(blockHash)
		if bEncoded == nil {
			return errors.New("the block is not existed")
		}
		result = base.DecodeBlock(bEncoded)
		return nil
	})
	return result, err
}

// GetNewestBlockHash 获取最新区块的hash
func (s *BoltStorage) GetNewestBlockHash() ([]byte, error) {
	var nbh []byte
	err := s.DB.View(func(tx *bolt.Tx) error {
		bhBucket := tx.Bucket([]byte(s.newestBlockHashBucket))
		// 该 bucket 中仅存储了这一个键值对
		nbh = bhBucket.Get([]byte(newestBlockHashKey))
		if nbh == nil {
			return errors.New("cannot find the newest block hash")
		}
		return nil
	})
	return nbh, err
}

func (s *BoltStorage) Close() error {
	return s.DB.Close()
}
