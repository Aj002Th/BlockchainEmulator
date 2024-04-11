package pbft

import (
	"encoding/json"
	"log"
)

// This module used in the blockChain using transaction relaying mechanism.
// "Raw" means that the pbft only make block consensus.
type RawRelayOutsideModule struct {
	pbftNode *PbftConsensusNode
}

// msgType canbe defined in message
func (rrom *RawRelayOutsideModule) HandleMessageOutsidePBFT(msgType MessageType, content []byte) bool {
	switch msgType {
	case CInject:
		rrom.handleInjectTx(content)
	default:
	}
	return true
}

func (rrom *RawRelayOutsideModule) handleInjectTx(content []byte) {
	it := new(InjectTxs)
	err := json.Unmarshal(content, it)
	if err != nil {
		log.Panic(err)
	}
	rrom.pbftNode.CurChain.TransactionPool.AddTransactionsToPool(it.Txs)
	rrom.pbftNode.pl.Printf("S%dN%d : has handled injected txs msg, txs: %d \n", rrom.pbftNode.ShardID, rrom.pbftNode.NodeID, len(it.Txs))
}
