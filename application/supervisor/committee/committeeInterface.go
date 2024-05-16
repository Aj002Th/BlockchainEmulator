package committee

import (
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
	"github.com/Aj002Th/BlockchainEmulator/data/base"
)

type CommitteeModule interface {
	txSending(txlist []*base.Transaction)
	HandleBlockInfo(*pbft.BlockInfoMsg)
	MsgSendingControl()
}
