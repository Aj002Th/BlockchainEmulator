package boot

import (
	"github.com/Aj002Th/BlockchainEmulator/params"
	"github.com/spf13/pflag"
)

type Args struct {
	nodeNum  int
	nodeID   int
	isClient bool
	frontend bool
}

// ParseAndBuildArg 旨在顺手把params里的值也设置了。也就是说params里的值是Preset。当然后期还要加一些逻辑。
func ParseAndBuildArg() Args {
	var a Args
	pflag.IntVarP(&a.nodeNum, "nodeNum", "N", 3, "indicate how many nodes of each shard are deployed")
	pflag.IntVarP(&a.nodeID, "nodeID", "n", 0, "id of this node, for example, 0")
	pflag.BoolVarP(&a.isClient, "client", "c", false, "whether this node is a client")
	pflag.BoolVarP(&a.frontend, "frontend", "f", false, "whether open web GUI monitor frontend")
	pflag.Parse()

	params.NodeNum = a.nodeNum

	return a
}
