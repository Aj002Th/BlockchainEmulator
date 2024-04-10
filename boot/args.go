package boot

import (
	"github.com/spf13/pflag"
)

type Args struct {
	shardNum int
	nodeNum  int
	shardID  int
	nodeID   int
	modID    int
	isClient bool
	isGen    bool
}

func ParseAndBuildArg() Args {
	var a Args
	pflag.IntVarP(&a.shardNum, "shardNum", "S", 2, "indicate that how many shards are deployed")
	pflag.IntVarP(&a.nodeNum, "nodeNum", "N", 4, "indicate how many nodes of each shard are deployed")
	pflag.IntVarP(&a.shardID, "shardID", "s", 0, "id of the shard to which this node belongs, for example, 0")
	pflag.IntVarP(&a.nodeID, "nodeID", "n", 0, "id of this node, for example, 0")
	pflag.IntVarP(&a.modID, "modID", "m", 3, "choice Committee Method,for example, 0, [CLPA_Broker,CLPA,Broker,Relay] ")
	pflag.BoolVarP(&a.isClient, "client", "c", false, "whether this node is a client")
	pflag.BoolVarP(&a.isGen, "gen", "g", false, "generation bat")
	pflag.Parse()
	return a
}
