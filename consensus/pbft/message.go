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

	CPbft   MessageType = "pbft"
	CInject MessageType = "inject"

	CBlockInfo MessageType = "BlockInfo"
	CSeqIDinfo MessageType = "SequenceID"
	CBooking   MessageType = "Booking"

	CKeepAlive MessageType = "KeepAlive"
)

var (
	BlockRequest RequestType = "Block"
)

type RawMessage struct {
	Content []byte
}

type Request struct {
	RequestType RequestType
	Msg         RawMessage
	ReqTime     time.Time
}

type PrePrepare struct {
	RequestMsg *Request
	Digest     []byte
	SeqID      uint64
}

type Prepare struct {
	Digest     []byte
	SeqID      uint64
	SenderNode *Node
}

type Commit struct {
	Digest     []byte
	SeqID      uint64
	SenderNode *Node
}

type Reply struct {
	MessageID  uint64
	SenderNode *Node
	Result     bool
}

type RequestOldMessage struct {
	SeqStartHeight uint64
	SeqEndHeight   uint64
	ServerNode     *Node
	SenderNode     *Node
}

type InjectTxs struct {
	Txs       []*base.Transaction
	ToShardID uint64
}

type BlockInfoMsg struct {
	BlockBodyLength int
	ExcutedTxs      []*base.Transaction
	Epoch           int

	ProposeTime   time.Time
	CommitTime    time.Time
	SenderShardID uint64

	TxpoolSize int
}

type BookingMsg struct {
	AvgCpuTime    float64 `json:"avgCpuTime"`
	DiskMetric    uint64  `json:"disk"`
	TotalUpload   int     `json:"tu"`
	TotalDownload int     `json:"td"`
	TotalTime     int64   `json:"tm"`
	NodeId        int     `json:"nodeid"`
}

type KeepAliveMsg struct {
	Msg string `json:"msg"`
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

	msgSend := MergeMessage(t, content)

	network.Tcp.Send(msgSend, addr)
}
