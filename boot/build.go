package boot

import (
	"log"
	"os/exec"
	"strconv"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor"
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/webapi"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
	"github.com/Aj002Th/BlockchainEmulator/data/chain"
	"github.com/Aj002Th/BlockchainEmulator/params"
)

// 初始化全局变量
func initGlobalConfig() {
	// 初始化 ip table
	if _, ok := params.IPmapNodeTable; !ok {
		log.Default().Println("提示：用户指定的IpTable不存在，按照约定生成默认的IpTable")
		params.IPmapNodeTable = make(map[uint64]string)

		for j := uint64(0); j < uint64(params.NodeNum); j++ {
			params.IPmapNodeTable[j] = "127.0.0.1:" + strconv.Itoa(10086+int(j))
		}
	} else {
		log.Default().Println("提示：发现用户指定的IpTable，将使用用户指定的IpTable")
	}
}

func makeChainConfig(nid uint64) *chain.Config {
	pcc := &chain.Config{
		NodeID:    nid,
		NodesNum:  uint64(params.NodeNum),
		BlockSize: uint64(params.MaxBlockSizeGlobal),
	}
	return pcc
}

func BuildSupervisor(self *App) {
	println("Build Sup")
	if self.args.frontend {
		webapi.GlobalProxy = webapi.NewGoodApiProxy()
		go webapi.RunApiServer()
		go webapi.RunFrontendServer()
		exec.Command("explorer.exe", "http://localhost:3000/monitor.html").Start() // 把浏览器拉起来
	} else {
		webapi.GlobalProxy = webapi.DumbProxy{}
	}

	webapi.GlobalProxy.Enqueue(webapi.Hello)

	sup := supervisor.NewSupervisor()
	sup.Run()
}

func BuildNewPbftNode(nid, nnm uint64) {
	worker := pbft.NewPbftNode(nid, makeChainConfig(nid))
	worker.Run()
}
