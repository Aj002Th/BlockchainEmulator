package pbft

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/data/base"
	"github.com/Aj002Th/BlockchainEmulator/data/chain"
	"github.com/Aj002Th/BlockchainEmulator/logger"
	"github.com/Aj002Th/BlockchainEmulator/params"
)

type PbftImplSimpleImpl struct {
	node *ConsensusNode
}

func (self *PbftImplSimpleImpl) doPropose() (bool, *Request) {

	block := self.node.CurChain.GenerateBlock()
	r := &Request{
		RequestType: BlockRequest,
		ReqTime:     time.Now(),
	}
	r.Msg.Content = block.Encode()

	return true, r
}

func (self *PbftImplSimpleImpl) doPreprepare(ppmsg *PrePrepare) bool {
	if self.node.CurChain.IsValidBlock(base.DecodeBlock(ppmsg.RequestMsg.Msg.Content)) != nil {
		self.node.pl.Printf("Node %d : not a valid block\n", self.node.NodeID)
		return false
	}
	self.node.pl.Printf("Node %d : the pre-prepare pbft is correct, putting it into the RequestPool. \n", self.node.NodeID)
	self.node.requestPool[string(ppmsg.Digest)] = ppmsg.RequestMsg

	return true
}

func (self *PbftImplSimpleImpl) doPrepare(pmsg *Prepare) bool {
	logger.Println("No operations are performed in Extra handle mod")
	return true
}

func PrintBlockChain(bc *chain.BlockChain) string {
	vals := []interface{}{
		bc.CurrentBlock.Header.Number,
		bc.CurrentBlock.Hash,
		bc.CurrentBlock.Header.StateRoot,
		bc.CurrentBlock.Header.Time,
		bc.BlockStorage,
	}
	res := fmt.Sprintf("%v\n", vals)
	logger.Println(res)
	return res
}

func (self *PbftImplSimpleImpl) doCommit(cmsg *Commit) bool {
	r := self.node.requestPool[string(cmsg.Digest)]

	block := base.DecodeBlock(r.Msg.Content)
	self.node.pl.Printf("Node %d : adding the block %d...now height = %d \n", self.node.NodeID, block.Header.Number, self.node.CurChain.CurrentBlock.Header.Number)
	self.node.CurChain.AddBlock(block)
	self.node.pl.Printf("Node %d : added the block %d... \n", self.node.NodeID, block.Header.Number)
	PrintBlockChain(self.node.CurChain)

	if self.node.NodeID == self.node.view {
		self.node.pl.Printf("Node %d : main node is trying to send pbft txs at height = %d \n", self.node.NodeID, block.Header.Number)

		txExcuted := make([]*base.Transaction, 0)
		for _, tx := range block.Body {
			txExcuted = append(txExcuted, tx)
		}

		bim := BlockInfoMsg{
			BlockBodyLength: len(block.Body),
			ExcutedTxs:      txExcuted,
			Epoch:           0,
			SenderShardID:   self.node.ShardID,
			ProposeTime:     r.ReqTime,
			CommitTime:      time.Now(),
			TxpoolSize:      len(self.node.CurChain.TransactionPool.Queue),
		}
		bByte, err := json.Marshal(bim)
		if err != nil {
			log.Panic()
		}
		MergeAndSend(CBlockInfo, bByte, params.SupervisorEndpoint, self.node.pl)
		self.node.pl.Printf("Node %d : sended excuted txs\n", self.node.NodeID)
	}
	return true
}
