package boot

import (
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor"
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/webapi"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
	"github.com/Aj002Th/BlockchainEmulator/data/chain"
	"github.com/Aj002Th/BlockchainEmulator/misc"
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

	// 根据环境变量为log生成prefix。如果没指定那就按分钟生成了。
	prefix := os.Getenv("BCEM_OUTPUT_PREFIX")
	if prefix == "" {
		dt := time.Now()
		prefix = dt.Format("2006-01-02T15:04:05")
	}
	var err error
	prefix, err = misc.CreateUniqueFolder(prefix)
	if err != nil {
		panic("unique folder create encountered an error.")
	}
	params.LogWrite_path = path.Join(params.LogWrite_path, prefix)
	params.DataWrite_path = path.Join(params.DataWrite_path, prefix)

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
		exec.Command("explorer.exe", "http://localhost:3000/monitor.html").Start() // 把浏览器拉起来
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

// 单机启动时的简便方法。
func StartNAtOnce(nnm uint64) {
	// 设置命令行参数的日志前缀
	dt := time.Now()
	prefix := dt.Format("2006-01-02T15:04:05") // 这是什么格式？我也不知道。
	os.Setenv("BCEM_OUTPUT_PREFIX", prefix)

	// 构造启动
	N := strconv.Itoa(int(nnm))
	// 依次启动各个
	for i := 0; i < int(nnm); i++ {
		n := strconv.Itoa(i)
		exec.Command("start", "cmd", "/k", "go", "run", "main.go", "-N", N, "-n", n).Start()
	}
	// 启动Supervisor
	exec.Command("start", "cmd", "/k", "go", "run", "main.go", "-N", N, "-c").Start()
}
