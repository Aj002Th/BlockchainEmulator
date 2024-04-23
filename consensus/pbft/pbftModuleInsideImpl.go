// addtional module for new consensus
package pbft

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/data/base"
	"github.com/Aj002Th/BlockchainEmulator/data/chain"
	"github.com/Aj002Th/BlockchainEmulator/params"
)

// simple implementation of pbftHandleModule interface ...
// only for block request and use transaction relay
type RawRelayPbftExtraHandleMod struct {
	node *PbftConsensusNode
	// pointer to pbft data
}

// propose request with different types
// 调用GenerateBlock产生一个区块。其实这个做的就是从交易池子里取出至多n个交易，并且弄成一个新区块。
// 更骚的是，没有新区块也会照样执行不误。。。所以Sup不给主节点东西他也会一直自娱自乐。
func (self *RawRelayPbftExtraHandleMod) HandleinPropose() (bool, *Request) {
	// new blocks
	block := self.node.CurChain.GenerateBlock()
	r := &Request{
		RequestType: BlockRequest,
		ReqTime:     time.Now(),
	}
	r.Msg.Content = block.Encode()

	return true, r
}

// the diy operation in preprepare
func (self *RawRelayPbftExtraHandleMod) HandleinPrePrepare(ppmsg *PrePrepare) bool {
	if self.node.CurChain.IsValidBlock(base.DecodeBlock(ppmsg.RequestMsg.Msg.Content)) != nil {
		self.node.pl.Printf("S%dN%d : not a valid block\n", self.node.ShardID, self.node.NodeID)
		return false
	}
	self.node.pl.Printf("S%dN%d : the pre-prepare pbft is correct, putting it into the RequestPool. \n", self.node.ShardID, self.node.NodeID)
	self.node.requestPool[string(ppmsg.Digest)] = ppmsg.RequestMsg
	// merge to be a prepare pbft
	return true
}

// the operation in prepare, and in pbft + tx relaying, this function does not need to do any.
func (self *RawRelayPbftExtraHandleMod) HandleinPrepare(pmsg *Prepare) bool {
	fmt.Println("No operations are performed in Extra handle mod")
	return true
}

// print the details of a blockchain
func PrintBlockChain(bc *chain.BlockChain) string {
	vals := []interface{}{
		bc.CurrentBlock.Header.Number,
		bc.CurrentBlock.Hash,
		bc.CurrentBlock.Header.StateRoot,
		bc.CurrentBlock.Header.Time,
		bc.BlockStorage,
	}
	res := fmt.Sprintf("%v\n", vals)
	fmt.Println(res)
	return res
}

// the operation in commit.
func (self *RawRelayPbftExtraHandleMod) HandleinCommit(cmsg *Commit) bool {
	r := self.node.requestPool[string(cmsg.Digest)]
	// requestType ...
	block := base.DecodeBlock(r.Msg.Content)
	self.node.pl.Printf("S%dN%d : adding the block %d...now height = %d \n", self.node.ShardID, self.node.NodeID, block.Header.Number, self.node.CurChain.CurrentBlock.Header.Number)
	self.node.CurChain.AddBlock(block)
	self.node.pl.Printf("S%dN%d : added the block %d... \n", self.node.ShardID, self.node.NodeID, block.Header.Number)
	PrintBlockChain(self.node.CurChain)

	// now try to relay txs to other shards (for main nodes)
	if self.node.NodeID == self.node.view {
		self.node.pl.Printf("S%dN%d : main node is trying to send relay txs at height = %d \n", self.node.ShardID, self.node.NodeID, block.Header.Number)
		// generate relay pool and collect txs excuted
		txExcuted := make([]*base.Transaction, 0)
		relay1Txs := make([]*base.Transaction, 0)
		for _, tx := range block.Body {
			txExcuted = append(txExcuted, tx)
		}
		// send txs excuted in this block to the listener
		// add more pbft to measure more metrics
		bim := BlockInfoMsg{
			BlockBodyLength: len(block.Body),
			ExcutedTxs:      txExcuted,
			Epoch:           0,
			Relay1Txs:       relay1Txs,
			Relay1TxNum:     uint64(len(relay1Txs)),
			SenderShardID:   self.node.ShardID,
			ProposeTime:     r.ReqTime,
			CommitTime:      time.Now(),
			TxpoolSize:      (len(self.node.CurChain.TransactionPool.Queue)),
		}
		bByte, err := json.Marshal(bim)
		if err != nil {
			log.Panic()
		}
		MergeAndSend(CBlockInfo, bByte, params.SupervisorEndpoint, self.node.pl)
		self.node.pl.Printf("S%dN%d : sended excuted txs\n", self.node.ShardID, self.node.NodeID)
		self.node.CurChain.TransactionPool.Locked()
		self.node.writeCSVline([]string{strconv.Itoa(len(self.node.CurChain.TransactionPool.Queue)), strconv.Itoa(len(txExcuted)), strconv.Itoa(int(bim.Relay1TxNum))})
		self.node.CurChain.TransactionPool.Unlocked()
	}
	return true
}
