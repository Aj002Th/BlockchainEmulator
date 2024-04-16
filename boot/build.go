package boot

import (
	"strconv"

	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
	"github.com/Aj002Th/BlockchainEmulator/data/chain"
	"github.com/Aj002Th/BlockchainEmulator/params"
)

// 等效于往一个全局的params里放一个params。这个params是大家共用的只读结构。最后返回一个莫名其妙的ChainCOnfig
func initConfig(nid uint64) *chain.Config {
	var i uint64 = 0
	if _, ok := params.IPmap_nodeTable[i]; !ok {
		params.IPmap_nodeTable[i] = make(map[uint64]string)
	}
	for j := uint64(0); j < uint64(params.NodeNum); j++ {
		params.IPmap_nodeTable[i][j] = "127.0.0.1:" + strconv.Itoa(28800+int(j)) // shard和node决定了ip
	}

	pcc := &chain.Config{
		NodeID:         nid,
		Nodes_perShard: uint64(params.NodeNum),
		BlockSize:      uint64(params.MaxBlockSize_global),
	}
	return pcc
}

func BuildNewPbftNode(nid, nnm uint64) {
	worker := pbft.NewPbftNode(nid, initConfig(nid), "Relay")
	if nid == 0 {
		worker.Run()
	} else {
		worker.Run()
	}
}
