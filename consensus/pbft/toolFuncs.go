package pbft

import (
	"crypto/sha256"
	"encoding/json"
	"log"
)

func (self *ConsensusNode) set2DMap(isPrePareConfirm bool, key string, val *Node) {
	if isPrePareConfirm {
		if _, ok := self.cntPrepareConfirm[key]; !ok {
			self.cntPrepareConfirm[key] = make(map[*Node]bool)
		}
		self.cntPrepareConfirm[key][val] = true
	} else {
		if _, ok := self.cntCommitConfirm[key]; !ok {
			self.cntCommitConfirm[key] = make(map[*Node]bool)
		}
		self.cntCommitConfirm[key][val] = true
	}
}

func (self *ConsensusNode) getNeighborNodes() []string {
	receiverNodes := make([]string, 0)
	for _, ip := range self.ipNodeTable[self.ShardID] {
		receiverNodes = append(receiverNodes, ip)
	}
	return receiverNodes
}

func getDigest(r *Request) []byte {
	b, err := json.Marshal(r)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}
