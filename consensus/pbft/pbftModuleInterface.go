package pbft

type ExtraOpInConsensus interface {
	// mining / message generation
	HandleinPropose() (bool, *Request)
	// checking
	HandleinPrePrepare(*PrePrepare) bool
	// nothing necessary
	HandleinPrepare(*Prepare) bool
	// confirming
	HandleinCommit(*Commit) bool
}

type OpInterShards interface {
	// operation inter-shards
	HandleMessageOutsidePBFT(MessageType, []byte) bool
}
