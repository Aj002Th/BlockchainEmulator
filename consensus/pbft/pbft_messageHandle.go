package pbft

import (
	"encoding/json"
	"log"

	"github.com/Aj002Th/BlockchainEmulator/logger"
)

func (self *ConsensusNode) handlePrePrepare(ppmsg *PrePrepare) {
	self.RunningNode.PrintNode()
	logger.Println("received the PrePrepare ...")
	// decode the message
	flag := false
	if digest := getDigest(ppmsg.RequestMsg); string(digest) != string(ppmsg.Digest) {
		self.pl.Printf("N%d : the digest is not consistent, so refuse to prepare. \n", self.NodeID)
	} else if self.sequenceID < ppmsg.SeqID {
		self.requestPool[string(getDigest(ppmsg.RequestMsg))] = ppmsg.RequestMsg
		self.height2Digest[ppmsg.SeqID] = string(getDigest(ppmsg.RequestMsg))
		self.pl.Printf("N%d : the Sequence id is not consistent, so refuse to prepare. \n", self.NodeID)
	} else {
		// do your operation in this interface
		flag = self.pbftImpl.doPreprepare(ppmsg)
		self.requestPool[string(getDigest(ppmsg.RequestMsg))] = ppmsg.RequestMsg
		self.height2Digest[ppmsg.SeqID] = string(getDigest(ppmsg.RequestMsg))
	}
	// if the message is true, broadcast the prepare message
	if flag {
		pre := Prepare{
			Digest:     ppmsg.Digest,
			SeqID:      ppmsg.SeqID,
			SenderNode: self.RunningNode,
		}
		prepareByte, err := json.Marshal(pre)
		if err != nil {
			log.Panic()
		}

		// broadcast
		MergeAndBroadcast(CPrepare, prepareByte, self.RunningNode.IPaddr, self.getNeighborNodes(), self.pl)
	}
}

func (self *ConsensusNode) handlePrepare(pmsg *Prepare) {
	self.pl.Printf("Node %d : received the Prepare ...\n", self.NodeID)

	if _, ok := self.requestPool[string(pmsg.Digest)]; !ok {
		self.pl.Printf("N%d : doesn't have the digest in the requst pool, refuse to commit\n", self.NodeID)
	} else if self.sequenceID < pmsg.SeqID {
		self.pl.Printf("N%d : inconsistent sequence ID, refuse to commit\n", self.NodeID)
	} else {
		// if needed more operations, implement interfaces
		self.pbftImpl.doPrepare(pmsg)

		self.set2DMap(true, string(pmsg.Digest), pmsg.SenderNode)
		cnt := 0
		for range self.cntPrepareConfirm[string(pmsg.Digest)] {
			cnt++
		}

		// the main node will not send the prepare message
		specifiedcnt := int(2 * self.maliciousNums)
		if self.NodeID != self.view {
			specifiedcnt -= 1
		}

		// if the node has received 2f messages (itself included), and it haven't committed, then it commit
		self.lock.Lock()
		defer self.lock.Unlock()
		if cnt >= specifiedcnt && !self.isCommitBroadcast[string(pmsg.Digest)] {
			self.pl.Printf("N%d : is going to commit\n", self.NodeID)
			// generate commit and broadcast
			c := Commit{
				Digest:     pmsg.Digest,
				SeqID:      pmsg.SeqID,
				SenderNode: self.RunningNode,
			}
			commitByte, err := json.Marshal(c)
			if err != nil {
				log.Panic()
			}
			self.isCommitBroadcast[string(pmsg.Digest)] = true
			MergeAndBroadcast(CCommit, commitByte, self.RunningNode.IPaddr, self.getNeighborNodes(), self.pl)
		}
	}
}

func (self *ConsensusNode) handleCommit(cmsg *Commit) {
	// decode the message
	self.pl.Printf("Node %d received the Commit from ...%d\n", self.NodeID, cmsg.SenderNode.NodeID)
	self.set2DMap(false, string(cmsg.Digest), cmsg.SenderNode)
	cnt := 0
	for range self.cntCommitConfirm[string(cmsg.Digest)] {
		cnt++
	}

	self.lock.Lock()
	defer self.lock.Unlock()
	// the main node will not send the prepare message
	required_cnt := int(2 * self.maliciousNums)
	if cnt >= required_cnt && !self.isReply[string(cmsg.Digest)] {
		self.pl.Printf("Node %d : has received 2f + 1 commits ... \n", self.NodeID)
		// if this node is left behind, so it need to requst blocks
		if _, ok := self.requestPool[string(cmsg.Digest)]; !ok {
			self.isReply[string(cmsg.Digest)] = true
			self.askForLock.Lock()
			// request the block
			sn := &Node{
				NodeID: self.view,
				IPaddr: self.nodeEndpointList[self.view],
			}
			orequest := RequestOldMessage{
				SeqStartHeight: self.sequenceID + 1,
				SeqEndHeight:   cmsg.SeqID,
				ServerNode:     sn,
				SenderNode:     self.RunningNode,
			}
			bromyte, err := json.Marshal(orequest)
			if err != nil {
				log.Panic()
			}

			self.pl.Printf("Node %d : is now requesting message (seq %d to %d) ... \n", self.NodeID, orequest.SeqStartHeight, orequest.SeqEndHeight)
			MergeAndSend(CRequestOldrequest, bromyte, orequest.ServerNode.IPaddr, self.pl)
		} else {
			// implement interface
			self.pbftImpl.doCommit(cmsg)
			self.isReply[string(cmsg.Digest)] = true
			self.pl.Printf("Node %d: this round of pbft %d is end \n", self.NodeID, self.sequenceID)
			self.sequenceID += 1
		}

		// if this node is a main node, then unlock the sequencelock
		if self.NodeID == self.view {
			self.sequenceLock.Unlock()
			self.pl.Printf("Node %d get sequenceLock unlocked...\n", self.NodeID)
		}
	}
}
