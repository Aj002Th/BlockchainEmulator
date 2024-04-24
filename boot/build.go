package boot

import (
	"os/exec"
	"strconv"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor"
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/webapi"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
	"github.com/Aj002Th/BlockchainEmulator/data/chain"
	"github.com/Aj002Th/BlockchainEmulator/params"
)

// 等效于往一个全局的params里放一个params。这个params是大家共用的只读结构。最后返回一个莫名其妙的ChainCOnfig
func initConfig() {
	if _, ok := params.IPmap_nodeTable[0]; !ok {
		params.IPmap_nodeTable[0] = make(map[uint64]string)
	}
	for j := uint64(0); j < uint64(params.NodeNum); j++ {
		params.IPmap_nodeTable[0][j] = "127.0.0.1:" + strconv.Itoa(28800+int(j)) // shard和node决定了ip
	}

}

func makeChainConfig(nid uint64) *chain.Config {

	pcc := &chain.Config{
		NodeID:         nid,
		Nodes_perShard: uint64(params.NodeNum),
		BlockSize:      uint64(params.MaxBlockSize_global),
	}
	return pcc
}

func BuildSupervisor(self *App) {
	println("Build Sup")
	if self.args.frontend {
		webapi.G_Proxy = webapi.NewGoodApiProxy()
		go webapi.RunApiServer()
		go webapi.RunFrontendServer()
		print("before exec start")
		exec.Command("start", "http://localhost:3000") // 把浏览器拉起来
		print("after exec start")
	} else {
		webapi.G_Proxy = webapi.DumbProxy{}
	}

	webapi.G_Proxy.Enqueue(webapi.Hello)

	sup := supervisor.NewSupervisor()
	time.Sleep(10000 * time.Millisecond) // TODO: 去掉丑陋的Sleep
	sup.Run()
}

func BuildNewPbftNode(nid, nnm uint64) {
	worker := pbft.NewPbftNode(nid, makeChainConfig(nid), "Relay")
	if nid == 0 {
		worker.Run()
	} else {
		worker.Run()
	}
}
