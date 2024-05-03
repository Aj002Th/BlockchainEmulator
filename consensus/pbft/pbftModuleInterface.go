package pbft

type ExtraOpInConsensus interface {
	// mining / message generation
	HandleInPropose() (bool, *Request)
	// checking
	HandleInPrePrepare(*PrePrepare) bool
	// nothing necessary
	HandleInPrepare(*Prepare) bool
	// confirming
	HandleInCommit(*Commit) bool
}

type OpInterShards interface {
	// operation inter-shards
	HandleMessageOutsidePBFT(MessageType, []byte) bool
}
