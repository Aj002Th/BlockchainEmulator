package chain

import (
	"errors"
	"log"
	"math/big"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/data/base"
	"github.com/Aj002Th/BlockchainEmulator/data/mpt"
	"github.com/Aj002Th/BlockchainEmulator/logger"
	"github.com/Aj002Th/BlockchainEmulator/storage/blockStorage"
	"github.com/Aj002Th/BlockchainEmulator/storage/stateStorage"
)

// BlockChain 区块链
type BlockChain struct {
	Config          *Config                   // 区块链配置
	CurrentBlock    *base.Block               // 最新的区块
	BlockStorage    blockStorage.BlockStorage // 区块持久化
	StateStorage    stateStorage.StateStorage // 状态树持久化
	Trie            *mpt.Trie                 // 状态树
	TransactionPool *base.TransactionPool     // 交易池
}

// GetTxTreeRoot 获取交易根，这个根可以用来检查区块中的交易
// 新建一个空的 mpt, 插入所有 transactions 并计算 root hash
func GetTxTreeRoot(txs []*base.Transaction) []byte {
	transactionTree := mpt.New(nil, stateStorage.NewMemKVStore())
	for _, tx := range txs {
		err := transactionTree.Put(tx.Hash, tx.Encode())
		if err != nil {
			log.Panic(err)
		}
	}
	return transactionTree.RootHash()
}

// SendTransactionsPool 将交易批量放入交易池
func (bc *BlockChain) SendTransactionsPool(txs []*base.Transaction) {
	bc.TransactionPool.AddTransactionsToPool(txs)
}

// GetUpdateStatusTrie 处理交易并更新状态树, 返回更新后的状态树根
func (bc *BlockChain) GetUpdateStatusTrie(txs []*base.Transaction) []byte {
	logger.Printf("The len of txs is %d\n", len(txs))

	// 如果没有交易产生, 说明状态树无需更新, 直接返回当前区块的状态树根
	if len(txs) == 0 {
		return bc.CurrentBlock.Header.StateRoot
	}

	// 处理交易, 此处做了简化, 平台没有实现用户签名验证
	st := bc.Trie
	cnt := 0
	for i, tx := range txs {
		// 更新 sender 账户(扣款)
		senderStateEncoded, _ := st.Get([]byte(tx.Sender))
		var senderState *base.AccountState
		if senderStateEncoded == nil {
			// logger.Println("missing account SENDER, now adding account")
			ib := new(big.Int)
			ib.Add(ib, InitBalance)
			senderState = &base.AccountState{
				Nonce:   uint64(i),
				Balance: ib,
			}
		} else {
			senderState = base.DecodeToAccountState(senderStateEncoded)
		}
		senderBalance := senderState.Balance
		if senderBalance.Cmp(tx.Value) == -1 {
			logger.Printf("the balance is less than the transfer amount\n")
			continue
		}
		senderState.Deduct(tx.Value)
		err := st.Put([]byte(tx.Sender), senderState.Encode())
		if err != nil {
			log.Panic(err)
		}
		cnt++

		// 更新 recipient 账户(取款)
		recipientStateEncoded, _ := st.Get([]byte(tx.Recipient))
		var recipientState *base.AccountState
		if recipientStateEncoded == nil {
			// logger.Println("missing account RECIPIENT, now adding account")
			ib := new(big.Int)
			ib.Add(ib, InitBalance)
			recipientState = &base.AccountState{
				Nonce:   uint64(i),
				Balance: ib,
			}
		} else {
			recipientState = base.DecodeToAccountState(recipientStateEncoded)
		}
		recipientState.Deposit(tx.Value)
		err = st.Put([]byte(tx.Recipient), recipientState.Encode())
		if err != nil {
			log.Panic(err)
		}
		cnt++
	}

	// 将状态树的修改持久化
	if cnt == 0 {
		return bc.CurrentBlock.Header.StateRoot
	}
	st.Commit()
	rt := st.RootHash()
	logger.Println("modified account number is ", cnt)
	return rt
}

// GenerateBlock 生成一个区块
func (bc *BlockChain) GenerateBlock() *base.Block {
	// 从交易池中取出交易
	txs := bc.TransactionPool.PackTransactions(bc.Config.BlockSize)
	bh := &base.BlockHeader{
		ParentBlockHash: bc.CurrentBlock.Hash,
		Number:          bc.CurrentBlock.Header.Number + 1,
		Time:            time.Now(),
	}
	// 处理交易并更新状态树
	rt := bc.GetUpdateStatusTrie(txs)

	bh.StateRoot = rt
	bh.TxRoot = GetTxTreeRoot(txs)
	b := base.NewBlock(bh, txs)
	b.Header.Miner = 0
	b.Hash = b.Header.Hash()
	return b
}

// NewGenisisBlock 生成创世区块, 每条区块链只会有一个创世区块
func (bc *BlockChain) NewGenisisBlock() *base.Block {
	body := make([]*base.Transaction, 0)
	bh := &base.BlockHeader{}
	b := base.NewBlock(bh, body)
	b.Hash = b.Header.Hash()
	return b
}

// AddGenisisBlock 将创世区块加入区块链
func (bc *BlockChain) AddGenisisBlock(gb *base.Block) {
	bc.BlockStorage.AddBlock(gb)
	newestHash, err := bc.BlockStorage.GetNewestBlockHash()
	if err != nil {
		log.Panic()
	}
	curb, err := bc.BlockStorage.GetBlock(newestHash)
	if err != nil {
		log.Panic()
	}
	bc.CurrentBlock = curb
}

// AddBlock 添加一个区块
func (bc *BlockChain) AddBlock(b *base.Block) {
	if b.Header.Number != bc.CurrentBlock.Header.Number+1 {
		logger.Println("the block height is not correct")
		return
	}
	// if this block is mined by the node, the transactions is no need to be handled again
	if b.Header.Miner != bc.Config.NodeID {
		rt := bc.GetUpdateStatusTrie(b.Body)
		logger.Println(bc.CurrentBlock.Header.Number+1, "the root = ", rt)
	}
	bc.CurrentBlock = b
	bc.BlockStorage.AddBlock(b)
}

// IsValidBlock 验证区块能否合法地上链
func (bc *BlockChain) IsValidBlock(b *base.Block) error {
	if string(b.Header.ParentBlockHash) != string(bc.CurrentBlock.Hash) {
		logger.Println("the parent block hash is not equal to the current block hash")
		return errors.New("the parent block hash is not equal to the current block hash")
	} else if string(GetTxTreeRoot(b.Body)) != string(b.Header.TxRoot) {
		logger.Println("the transaction root is wrong")
		return errors.New("the transaction root is wrong")
	}
	return nil
}

// AddAccounts 新建账户
func (bc *BlockChain) AddAccounts(ac []string, as []*base.AccountState) {
	logger.Printf("The len of accounts is %d, now adding the accounts\n", len(ac))

	bh := &base.BlockHeader{
		ParentBlockHash: bc.CurrentBlock.Hash,
		Number:          bc.CurrentBlock.Header.Number + 1,
		Time:            time.Time{},
	}

	// 将账户信息插入状态树, 并持久化
	rt := bc.CurrentBlock.Header.StateRoot
	if len(ac) != 0 {
		st := bc.Trie
		for i, addr := range ac {
			ib := new(big.Int)
			ib.Add(ib, as[i].Balance)
			newState := &base.AccountState{
				Balance: ib,
				Nonce:   as[i].Nonce,
			}
			err := st.Put([]byte(addr), newState.Encode())
			if err != nil {
				log.Panic(err)
			}
		}
		st.Commit()
		rt = st.RootHash()
	}

	// 创建对应的区块
	emptyTxs := make([]*base.Transaction, 0)
	bh.StateRoot = rt
	bh.TxRoot = GetTxTreeRoot(emptyTxs)
	b := base.NewBlock(bh, emptyTxs)
	b.Header.Miner = 0
	b.Hash = b.Header.Hash()

	// 将区块上链
	bc.CurrentBlock = b
	bc.BlockStorage.AddBlock(b)
}

// FetchAccounts 获取账户信息
func (bc *BlockChain) FetchAccounts(addrs []string) []*base.AccountState {
	result := make([]*base.AccountState, 0)
	st := bc.Trie
	for _, addr := range addrs {
		asEncoded, _ := st.Get([]byte(addr))

		// 如果账户不存在，则返回一个初始化的账户信息
		// 否则返回从状态树中取出的账户信息
		var accountState *base.AccountState
		if asEncoded == nil {
			ib := new(big.Int)
			ib.Add(ib, InitBalance)
			accountState = &base.AccountState{
				Nonce:   uint64(0),
				Balance: ib,
			}
		} else {
			accountState = base.DecodeToAccountState(asEncoded)
		}

		result = append(result, accountState)
	}
	return result
}

// CloseBlockChain 关闭区块链, 释放资源
func (bc *BlockChain) CloseBlockChain() {
	_ = bc.BlockStorage.Close()
}

// NewBlockChain 创建区块链
func NewBlockChain(conf *Config, stateDB stateStorage.StateStorage, blockDB blockStorage.BlockStorage) (*BlockChain, error) {
	logger.Println("Generating a new blockchain")
	bc := &BlockChain{
		Config:          conf,
		CurrentBlock:    nil,
		BlockStorage:    blockDB,
		StateStorage:    stateDB,
		Trie:            mpt.New(nil, stateDB),
		TransactionPool: base.NewTransactionPool(),
	}

	// 如果区块链是空的, 则创建创世区块
	curHash, err := bc.BlockStorage.GetNewestBlockHash()
	if err != nil {
		logger.Println("Get newest block hash err")
		// 如果是没能找到最新区块的hash，那么说明这个存储里面是空的，需要创建创世区块
		if err.Error() == "cannot find the newest block hash" {
			genisisBlock := bc.NewGenisisBlock()
			bc.AddGenisisBlock(genisisBlock)
			logger.Println("New genisis block")
			return bc, nil
		}
		log.Panic(err)
	}

	// 获取最新的区块作为链尾
	logger.Println("Existing blockchain found")
	curb, err := bc.BlockStorage.GetBlock(curHash)
	if err != nil {
		log.Panic(err)
	}
	bc.CurrentBlock = curb

	logger.Println("The status trie can be built")
	logger.Println("Generated a new blockchain successfully")
	return bc, nil
}
