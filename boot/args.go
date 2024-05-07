package boot

import (
	"fmt"
	"os"
	"path"

	"github.com/Aj002Th/BlockchainEmulator/params"
	"github.com/spf13/pflag"
)

type Args struct {
	nodeID        int
	isClient      bool
	frontend      bool
	interpretArgs bool
}

func PrintCurrentArgs() {
	fmt.Printf("The params-set will be:\n")
	fmt.Printf("params.NodeNum               = %v\n", params.NodeNum)
	fmt.Printf("params.BlockInterval         = %v\n", params.BlockInterval)
	fmt.Printf("params.MaxBlockSizeGlobal    = %v\n", params.MaxBlockSizeGlobal)
	fmt.Printf("params.InjectSpeed           = %v\n", params.InjectSpeed)
	fmt.Printf("params.TotalDataSize         = %v\n", params.TotalDataSize)
	fmt.Printf("params.BatchSize             = %v\n", params.BatchSize)
	fmt.Printf("params.LogWritePath          = %v\n", params.LogWritePath)
	fmt.Printf("params.DataWritePath         = %v\n", params.DataWritePath)
	fmt.Printf("params.RecordWritePath       = %v\n", params.RecordWritePath)
	fmt.Printf("params.SupervisorEndpoint    = %v\n", params.SupervisorEndpoint)
	fmt.Printf("params.FileInput             = %v\n", params.FileInput)
}

// ParseAndBuildArg 旨在顺手把params里的值也设置了
func ParseAndBuildArg() Args {
	// 启动控制
	var a Args
	pflag.IntVarP(&a.nodeID, "nodeID", "n", 0, "id of this node, for example, 0")
	pflag.BoolVarP(&a.isClient, "client", "c", false, "whether this node is a client")
	pflag.BoolVarP(&a.frontend, "frontend", "f", false, "whether open web GUI monitor frontend")
	pflag.BoolVarP(&a.interpretArgs, "interpret", "", false, "whether interpret the currently set args")

	// 参数注入
	// todo: add shorthand
	pflag.IntVarP(&params.NodeNum, "NodeNum", "N", 3, "indicate how many nodes of each shard are deployed")
	pflag.IntVarP(&params.BlockInterval, "BlockInterval", "", 5000, "generate new block interval")
	pflag.IntVarP(&params.MaxBlockSizeGlobal, "MaxBlockSizeGlobal", "", 2000, "the block contains the maximum number of transactions")
	pflag.IntVarP(&params.InjectSpeed, "InjectSpeed", "", 2000, "the transaction inject speed")
	pflag.IntVarP(&params.TotalDataSize, "TotalDataSize", "", 150000, "the total number of txs")
	pflag.IntVarP(&params.BatchSize, "BatchSize", "", 15000, "supervisor read a batch of txs then send them, it should be larger than inject speed")
	pflag.StringVarP(&params.LogWritePath, "LogWritePath", "", "./log", "log output path")
	pflag.StringVarP(&params.DataWritePath, "DataWritePath", "", "./result", "measurement data result output path")
	pflag.StringVarP(&params.RecordWritePath, "RecordWritePath", "", "./record", "record output path")
	pflag.StringVarP(&params.SupervisorEndpoint, "SupervisorEndpoint", "", "127.0.0.1:18800", "supervisor ip address")
	pflag.StringVarP(&params.FileInput, "FileInput", "i", "./BlockTransaction.csv", "the raw BlockTransaction data path")

	pflag.Parse()

	prefix := os.Getenv("BCEM_OUTPUT_PREFIX")
	if prefix == "" {
		panic("Set the BCEM_OUTPUT_PREFIX env var!")
	}
	// 生成结果文件对应的输出目录
	// prefix := time.Now().Format("01-02-2006-15-04-05")
	lPath := path.Join(params.LogWritePath, prefix)
	dPath := path.Join(params.DataWritePath, prefix)
	rPath := path.Join(params.RecordWritePath, prefix)
	// 覆写全局变量
	params.LogWritePath = lPath
	params.DataWritePath = dPath
	params.RecordWritePath = rPath

	return a
}
