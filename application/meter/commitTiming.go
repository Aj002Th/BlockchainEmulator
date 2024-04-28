package meter

import "github.com/Aj002Th/BlockchainEmulator/consensus/pbft"

// 包括TPS，TCL那些

type CommitCtx struct {
}

// 在sup端计算的时候用。
func Feed(commitCtx *CommitCtx, bim pbft.BlockInfoMsg) {

}
