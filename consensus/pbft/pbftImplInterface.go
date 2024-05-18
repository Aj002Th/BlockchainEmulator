package pbft

type PbftImplInterface interface {
	doPropose() (bool, *Request)
	doPreprepare(*PrePrepare) bool
	doPrepare(*Prepare) bool
	doCommit(*Commit) bool
}
