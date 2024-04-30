package pbft

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/Aj002Th/BlockchainEmulator/params"
)

// set 2d map, only for pbft maps, if the first parameter is true, then set the cntPrepareConfirm map,
// otherwise, cntCommitConfirm map will be set
func (self *PbftConsensusNode) set2DMap(isPrePareConfirm bool, key string, val *Node) {
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

// get neighbor nodes in a shard
func (self *PbftConsensusNode) getNeighborNodes() []string {
	receiverNodes := make([]string, 0)
	for _, ip := range self.ip_nodeTable[self.ShardID] {
		receiverNodes = append(receiverNodes, ip)
	}
	return receiverNodes
}

func (self *PbftConsensusNode) writeCSVline(str []string) {
	dirpath := path.Join(params.DataWrite_path, "pbft_"+strconv.Itoa(int(1)))
	err := os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}

	targetPath := dirpath + "/Shard" + strconv.Itoa(int(self.ShardID)) + strconv.Itoa(int(1)) + ".csv"
	f, err := os.Open(targetPath)
	if err != nil && os.IsNotExist(err) {
		file, er := os.Create(targetPath)
		if er != nil {
			panic(er)
		}
		defer file.Close()

		w := csv.NewWriter(file)
		title := []string{"txpool size", "tx", "ctx"}
		w.Write(title)
		w.Flush()
		w.Write(str)
		w.Flush()
	} else {
		file, err := os.OpenFile(targetPath, os.O_APPEND|os.O_RDWR, 0666)

		if err != nil {
			log.Panic(err)
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		err = writer.Write(str)
		if err != nil {
			log.Panic()
		}
		writer.Flush()
	}

	f.Close()
}

// get the digest of request
func getDigest(r *Request) []byte {
	b, err := json.Marshal(r)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}
