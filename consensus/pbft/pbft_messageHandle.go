package pbft

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/network"
	"github.com/Aj002Th/BlockchainEmulator/params"
)

// this func is only invoked by main node
func (p *PbftConsensusNode) Propose() {
	if p.view != p.NodeID {
		return
	}
	for {
		select {
		case <-p.pStop:
			p.pl.Printf("S%dN%d stop...\n", p.ShardID, p.NodeID)
			return
		default:
		}
		time.Sleep(time.Duration(int64(params.Block_Interval)) * time.Millisecond)

		p.sequenceLock.Lock()
		p.pl.Printf("S%dN%d get sequenceLock locked, now trying to propose...\n", p.ShardID, p.NodeID)
		// propose
		// implement interface to generate propose
		_, r := p.ihm.HandleinPropose()

		digest := getDigest(r)
		p.requestPool[string(digest)] = r
		p.pl.Printf("S%dN%d put the request into the pool ...\n", p.ShardID, p.NodeID)

		ppmsg := PrePrepare{
			RequestMsg: r,
			Digest:     digest,
			SeqID:      p.sequenceID,
		}
		p.height2Digest[p.sequenceID] = string(digest)
		// marshal and broadcast
		ppbyte, err := json.Marshal(ppmsg)
		if err != nil {
			log.Panic()
		}
		msg_send := MergeMessage(CPrePrepare, ppbyte)
		network.Tcp.Broadcast(p.RunningNode.IPaddr, p.getNeighborNodes(), msg_send)
	}
}

func (p *PbftConsensusNode) handlePrePrepare(content []byte) {
	p.RunningNode.PrintNode()
	fmt.Println("received the PrePrepare ...")
	// decode the message
	ppmsg := new(PrePrepare)
	err := json.Unmarshal(content, ppmsg)
	if err != nil {
		log.Panic(err)
	}
	flag := false
	if digest := getDigest(ppmsg.RequestMsg); string(digest) != string(ppmsg.Digest) {
		p.pl.Printf("S%dN%d : the digest is not consistent, so refuse to prepare. \n", p.ShardID, p.NodeID)
	} else if p.sequenceID < ppmsg.SeqID {
		p.requestPool[string(getDigest(ppmsg.RequestMsg))] = ppmsg.RequestMsg
		p.height2Digest[ppmsg.SeqID] = string(getDigest(ppmsg.RequestMsg))
		p.pl.Printf("S%dN%d : the Sequence id is not consistent, so refuse to prepare. \n", p.ShardID, p.NodeID)
	} else {
		// do your operation in this interface
		flag = p.ihm.HandleinPrePrepare(ppmsg)
		p.requestPool[string(getDigest(ppmsg.RequestMsg))] = ppmsg.RequestMsg
		p.height2Digest[ppmsg.SeqID] = string(getDigest(ppmsg.RequestMsg))
	}
	// if the message is true, broadcast the prepare message
	if flag {
		pre := Prepare{
			Digest:     ppmsg.Digest,
			SeqID:      ppmsg.SeqID,
			SenderNode: p.RunningNode,
		}
		prepareByte, err := json.Marshal(pre)
		if err != nil {
			log.Panic()
		}
		// broadcast
		msg_send := MergeMessage(CPrepare, prepareByte)
		network.Tcp.Broadcast(p.RunningNode.IPaddr, p.getNeighborNodes(), msg_send)
		p.pl.Printf("S%dN%d : has broadcast the prepare message \n", p.ShardID, p.NodeID)
	}
}

func (p *PbftConsensusNode) handlePrepare(content []byte) {
	p.pl.Printf("S%dN%d : received the Prepare ...\n", p.ShardID, p.NodeID)
	// decode the message
	pmsg := new(Prepare)
	err := json.Unmarshal(content, pmsg)
	if err != nil {
		log.Panic(err)
	}

	if _, ok := p.requestPool[string(pmsg.Digest)]; !ok {
		p.pl.Printf("S%dN%d : doesn't have the digest in the requst pool, refuse to commit\n", p.ShardID, p.NodeID)
	} else if p.sequenceID < pmsg.SeqID {
		p.pl.Printf("S%dN%d : inconsistent sequence ID, refuse to commit\n", p.ShardID, p.NodeID)
	} else {
		// if needed more operations, implement interfaces
		p.ihm.HandleinPrepare(pmsg)

		p.set2DMap(true, string(pmsg.Digest), pmsg.SenderNode)
		cnt := 0
		for range p.cntPrepareConfirm[string(pmsg.Digest)] {
			cnt++
		}
		// the main node will not send the prepare message
		specifiedcnt := int(2 * p.malicious_nums)
		if p.NodeID != p.view {
			specifiedcnt -= 1
		}

		// if the node has received 2f messages (itself included), and it haven't committed, then it commit
		p.lock.Lock()
		defer p.lock.Unlock()
		if cnt >= specifiedcnt && !p.isCommitBordcast[string(pmsg.Digest)] {
			p.pl.Printf("S%dN%d : is going to commit\n", p.ShardID, p.NodeID)
			// generate commit and broadcast
			c := Commit{
				Digest:     pmsg.Digest,
				SeqID:      pmsg.SeqID,
				SenderNode: p.RunningNode,
			}
			commitByte, err := json.Marshal(c)
			if err != nil {
				log.Panic()
			}
			msg_send := MergeMessage(CCommit, commitByte)
			network.Tcp.Broadcast(p.RunningNode.IPaddr, p.getNeighborNodes(), msg_send)
			p.isCommitBordcast[string(pmsg.Digest)] = true
			p.pl.Printf("S%dN%d : commit is broadcast\n", p.ShardID, p.NodeID)
		}
	}
}

func (p *PbftConsensusNode) handleCommit(content []byte) {
	// decode the message
	cmsg := new(Commit)
	err := json.Unmarshal(content, cmsg)
	if err != nil {
		log.Panic(err)
	}
	p.pl.Printf("S%dN%d received the Commit from ...%d\n", p.ShardID, p.NodeID, cmsg.SenderNode.NodeID)
	p.set2DMap(false, string(cmsg.Digest), cmsg.SenderNode)
	cnt := 0
	for range p.cntCommitConfirm[string(cmsg.Digest)] {
		cnt++
	}

	p.lock.Lock()
	defer p.lock.Unlock()
	// the main node will not send the prepare message
	required_cnt := int(2 * p.malicious_nums)
	if cnt >= required_cnt && !p.isReply[string(cmsg.Digest)] {
		p.pl.Printf("S%dN%d : has received 2f + 1 commits ... \n", p.ShardID, p.NodeID)
		// if this node is left behind, so it need to requst blocks
		if _, ok := p.requestPool[string(cmsg.Digest)]; !ok {
			p.isReply[string(cmsg.Digest)] = true
			p.askForLock.Lock()
			// request the block
			sn := &Node{
				NodeID:  p.view,
				ShardID: p.ShardID,
				IPaddr:  p.ip_nodeTable[p.ShardID][p.view],
			}
			orequest := RequestOldMessage{
				SeqStartHeight: p.sequenceID + 1,
				SeqEndHeight:   cmsg.SeqID,
				ServerNode:     sn,
				SenderNode:     p.RunningNode,
			}
			bromyte, err := json.Marshal(orequest)
			if err != nil {
				log.Panic()
			}

			p.pl.Printf("S%dN%d : is now requesting message (seq %d to %d) ... \n", p.ShardID, p.NodeID, orequest.SeqStartHeight, orequest.SeqEndHeight)
			msg_send := MergeMessage(CRequestOldrequest, bromyte)
			network.Tcp.Send(msg_send, orequest.ServerNode.IPaddr)
		} else {
			// implement interface
			p.ihm.HandleinCommit(cmsg)
			p.isReply[string(cmsg.Digest)] = true
			p.pl.Printf("S%dN%d: this round of pbft %d is end \n", p.ShardID, p.NodeID, p.sequenceID)
			p.sequenceID += 1
		}

		// if this node is a main node, then unlock the sequencelock
		if p.NodeID == p.view {
			p.sequenceLock.Unlock()
			p.pl.Printf("S%dN%d get sequenceLock unlocked...\n", p.ShardID, p.NodeID)
		}
	}
}
