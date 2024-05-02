// definition of node and shard

package pbft

import "github.com/Aj002Th/BlockchainEmulator/logger"

type Node struct {
	NodeID  uint64
	ShardID uint64
	IPaddr  string
}

func (n *Node) PrintNode() {
	v := []interface{}{
		n.NodeID,
		n.ShardID,
		n.IPaddr,
	}
	logger.Printf("%v\n", v)
}
