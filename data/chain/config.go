package chain

import "math/big"

// Config 区块链配置
type Config struct {
	NodeID         uint64
	BlockSize      uint64
	Nodes_perShard uint64
	//BlockInterval uint64
	//InjectSpeed   uint64
}

var (
	// 账户的初始余额
	InitBalance, _ = new(big.Int).SetString("100000000000000000000000000000000000000000000", 10)
)
