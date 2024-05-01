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
	"github.com/Aj002Th/BlockchainEmulator/params"
)

// 初始化params。原params内容只是Preset，现在覆写它。包括Endpoint列表、输出路径。
func initConfig() {
	if _, ok := params.IPmapNodeTable[0]; !ok {
		params.IPmapNodeTable[0] = make(map[uint64]string)
	}
	for j := uint64(0); j < uint64(params.NodeNum); j++ {
		params.IPmapNodeTable[0][j] = "127.0.0.1:" + strconv.Itoa(28800+int(j)) // shard和node决定了ip
	}

	prefix := os.Getenv("BCEM_OUTPUT_PREFIX")
	if prefix == "" {
		panic("Set the BCEM_OUTPUT_PREFIX env var!")
	}

	// 暂定两个文件夹。然后试着生成。不行就加后缀。
	lPath := path.Join(params.LogWritePath, prefix)
	dPath := path.Join(params.DataWritePath, prefix)
	rPath := path.Join(params.RecordWritePath, prefix)

	// 覆写全局变量
	params.LogWritePath = lPath
	params.DataWritePath = dPath
	params.RecordWritePath = rPath
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
	panic("Don't you use this -g flag!")
	// // 设置命令行参数的前缀的环境变量。没有指定就自己生成。
	// prefix := os.Getenv("BCEM_OUTPUT_PREFIX")
	// if prefix == "" {
	// 	dt := time.Now()
	// 	prefix = dt.Format("BCEM-20060102-150405")
	// }
	// os.Setenv("BCEM_OUTPUT_PREFIX", prefix)

	// // 获取项目根目录的正式路径。
	// cwd, err := os.Getwd()
	// if err != nil {
	// 	panic("get cwd error.")
	// }

	// 构造启动
	N := strconv.Itoa(int(nnm))
	// 依次启动各个
	for i := 0; i < int(nnm); i++ {
		n := strconv.Itoa(i)
		cmd := exec.Command("cmd", "/k", "start", "go", "run", "main.go", "-N", N, "-n", n)
		// cmd.Dir = cwd
		cmd.Start()
	}
	// 启动Supervisor
	cmd := exec.Command("cmd", "/k", "start", "go", "run", "main.go", "-N", N, "-c")
	// cmd.Dir = cwd
	cmd.Start()
}
