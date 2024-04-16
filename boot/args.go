package boot

import (
	"github.com/spf13/pflag"
)

type Args struct {
	nodeNum  int
	nodeID   int
	isClient bool
	frontend bool
}

func ParseAndBuildArg() Args {
	var a Args
	pflag.IntVarP(&a.nodeNum, "nodeNum", "N", 3, "indicate how many nodes of each shard are deployed")
	pflag.IntVarP(&a.nodeID, "nodeID", "n", 0, "id of this node, for example, 0")
	pflag.BoolVarP(&a.isClient, "client", "c", false, "whether this node is a client")
	pflag.BoolVarP(&a.frontend, "frontend", "f", false, "whether open web GUI monitor frontend")
	pflag.Parse()
	return a
}
