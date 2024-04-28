package meter

import "github.com/Aj002Th/BlockchainEmulator/signal"

var txCount uint64
var BlockCount uint64

// 数数，很简单的。在节点上计算就好了。
func onCommited(txs int) {
	txCount += uint64(txs)
	BlockCount += 1
}

func StartCnt() {
	sig := signal.GetSignalByName[int]("PbftOnCommitTxs")
	sig.Connect(onCommited)
}
