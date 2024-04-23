package pbft

import (
	"encoding/json"
	"log"
)

// This module used in the blockChain using transaction relaying mechanism.
// "Raw" means that the pbft only make block consensus.
type RawRelayOutsideModule struct {
	node *PbftConsensusNode
}

// msgType canbe defined in message
func (self *RawRelayOutsideModule) HandleMessageOutsidePBFT(msgType MessageType, content []byte) bool {
	switch msgType {
	case CInject:
		self.handleInjectTx(content)
	default:
	}
	return true
}

func (self *RawRelayOutsideModule) handleInjectTx(content []byte) {
	it := new(InjectTxs)
	err := json.Unmarshal(content, it)
	if err != nil {
		log.Panic(err)
	}
	self.node.CurChain.TransactionPool.AddTransactionsToPool(it.Txs)
	self.node.pl.Printf("S%dN%d : has handled injected txs msg, txs: %d \n", self.node.ShardID, self.node.NodeID, len(it.Txs))
}
