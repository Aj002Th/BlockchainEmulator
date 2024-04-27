package committee

import "github.com/Aj002Th/BlockchainEmulator/consensus/pbft"

type CommitteeModule interface {
	HandleBlockInfo(*pbft.BlockInfoMsg)
	MsgSendingControl()
}
