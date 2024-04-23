package pbft

import (
	"log"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/data/base"
	"github.com/Aj002Th/BlockchainEmulator/network"
)

var prefixMSGtypeLen = 30

type MessageType string
type RequestType string

const (
	CPrePrepare        MessageType = "preprepare"
	CPrepare           MessageType = "prepare"
	CCommit            MessageType = "commit"
	CRequestOldrequest MessageType = "requestOldrequest"
	CStop              MessageType = "stop"

	CRelay  MessageType = "relay"
	CInject MessageType = "inject"

	CBlockInfo MessageType = "BlockInfo"
	CSeqIDinfo MessageType = "SequenceID"
)

var (
	BlockRequest RequestType = "Block"
	// add more types
	// ...
)

type RawMessage struct {
	Content []byte // the content of raw message, txs and blocks (most cases) included
}

type Request struct {
	RequestType RequestType
	Msg         RawMessage // request message
	ReqTime     time.Time  // request time
}

type PrePrepare struct {
	RequestMsg *Request // the request message should be pre-prepared
	Digest     []byte   // the digest of this request, which is the only identifier
	SeqID      uint64
}

type Prepare struct {
	Digest     []byte // To identify which request is prepared by this node
	SeqID      uint64
	SenderNode *Node // To identify who send this message
}

type Commit struct {
	Digest     []byte // To identify which request is prepared by this node
	SeqID      uint64
	SenderNode *Node // To identify who send this message
}

type Reply struct {
	MessageID  uint64
	SenderNode *Node
	Result     bool
}

type RequestOldMessage struct {
	SeqStartHeight uint64
	SeqEndHeight   uint64
	ServerNode     *Node // send this request to the server node
	SenderNode     *Node
}

type SendOldMessage struct {
	SeqStartHeight uint64
	SeqEndHeight   uint64
	OldRequest     []*Request
	SenderNode     *Node
}

type InjectTxs struct {
	Txs       []*base.Transaction
	ToShardID uint64
}

type BlockInfoMsg struct {
	BlockBodyLength int
	ExcutedTxs      []*base.Transaction // txs which are excuted completely
	Epoch           int

	ProposeTime   time.Time // record the propose time of this block (txs)
	CommitTime    time.Time // record the commit time of this block (txs)
	SenderShardID uint64

	// for transaction relay
	Relay1TxNum uint64              // the number of cross shard txs
	Relay1Txs   []*base.Transaction // cross transactions in chain first time

	TxpoolSize int
}

type SeqIDinfo struct {
	SenderShardID uint64
	SenderSeq     uint64
}

func MergeMessage(msgType MessageType, content []byte) []byte {
	b := make([]byte, prefixMSGtypeLen)
	for i, v := range []byte(msgType) {
		b[i] = v
	}
	merge := append(b, content...)
	return merge
}

func SplitMessage(message []byte) (MessageType, []byte) {
	msgTypeBytes := message[:prefixMSGtypeLen]
	msgType_pruned := make([]byte, 0)
	for _, v := range msgTypeBytes {
		if v != byte(0) {
			msgType_pruned = append(msgType_pruned, v)
		}
	}
	msgType := string(msgType_pruned)
	content := message[prefixMSGtypeLen:]
	return MessageType(msgType), content
}

func MergeAndSend(t MessageType, content []byte, addr string, logger *log.Logger) {
	if len(content) > 2000 {
		logger.Printf("Sending a %v: %v\n", t, string(content[:2000]))
	} else {
		logger.Printf("Sending a %v: %v\n", t, string(content))
	}

	msg_send := MergeMessage(t, content)

	network.Tcp.Send(msg_send, addr)
}
