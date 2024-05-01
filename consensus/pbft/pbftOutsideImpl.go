package pbft

import (
	"encoding/json"
	"log"
)

// RawPbftOutsideModule
// This module used in the blockChain using transaction pbfting mechanism.
// "Raw" means that the pbft only make block consensus.
type RawPbftOutsideModule struct {
	node *PbftConsensusNode
}

// HandleMessageOutsidePBFT msgType canbe defined in message
func (self *RawPbftOutsideModule) HandleMessageOutsidePBFT(msgType MessageType, content []byte) bool {
	switch msgType {
	case CInject:
		self.handleInjectTx(content)
	default:
	}
	return true
}

func (self *RawPbftOutsideModule) handleInjectTx(content []byte) {
	it := new(InjectTxs)
	err := json.Unmarshal(content, it)
	if err != nil {
		log.Panic(err)
	}
	self.node.CurChain.TransactionPool.AddTransactionsToPool(it.Txs)
	self.node.pl.Printf("S%dN%d : has handled injected txs msg, txs: %d \n", self.node.ShardID, self.node.NodeID, len(it.Txs))
}
