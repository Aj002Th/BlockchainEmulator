package base

import (
	"sync"
	"time"
)

// TransactionPool 交易池
type TransactionPool struct {
	Queue []*Transaction // 交易队列
	lock  sync.Mutex     // 交易池锁, 保证并发安全
}

// NewTransactionPool 创建交易池
func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		Queue: make([]*Transaction, 0),
	}
}

// AddTransactionToPool 添加交易
func (txpool *TransactionPool) AddTransactionToPool(tx *Transaction) {
	txpool.lock.Lock()
	defer txpool.lock.Unlock()
	if tx.Time.IsZero() {
		tx.Time = time.Now()
	}
	txpool.Queue = append(txpool.Queue, tx)
}

// AddTransactionsToPool 批量添加交易
func (txpool *TransactionPool) AddTransactionsToPool(txs []*Transaction) {
	txpool.lock.Lock()
	defer txpool.lock.Unlock()
	for _, tx := range txs {
		if tx.Time.IsZero() {
			tx.Time = time.Now()
		}
		txpool.Queue = append(txpool.Queue, tx)
	}
}

// PackTransactions 获取交易
// maxTxs: 能获取到的最大交易数目, 当前交易池中交易数目小于maxTxs时, 返回全部交易, 否则返回maxTxs个交易
func (txpool *TransactionPool) PackTransactions(maxTxs uint64) []*Transaction {
	txpool.lock.Lock()
	defer txpool.lock.Unlock()
	txNum := maxTxs
	if uint64(len(txpool.Queue)) < txNum {
		txNum = uint64(len(txpool.Queue))
	}
	txsPacked := txpool.Queue[:txNum]
	txpool.Queue = txpool.Queue[txNum:]
	return txsPacked
}

// Locked 获取锁
func (txpool *TransactionPool) Locked() {
	txpool.lock.Lock()
}

// Unlocked 释放锁
func (txpool *TransactionPool) Unlocked() {
	txpool.lock.Unlock()
}

// Size 获取交易池中的交易数目
func (txpool *TransactionPool) Size() int {
	txpool.lock.Lock()
	defer txpool.lock.Unlock()
	return len(txpool.Queue)
}

// ClearQueue 清空交易池
func (txpool *TransactionPool) ClearQueue() {
	txpool.lock.Lock()
	defer txpool.lock.Unlock()
	txpool.Queue = make([]*Transaction, 0)
}
