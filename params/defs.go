package params

type ChainConfig struct {
	ChainID        uint64
	NodeID         uint64
	ShardID        uint64
	Nodes_perShard uint64
	ShardNums      uint64
	BlockSize      uint64
	BlockInterval  uint64
	InjectSpeed    uint64

	// used in transaction relaying, useless in brokerchain mechanism
	MaxRelayBlockSize uint64
}
