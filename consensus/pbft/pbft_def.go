package pbft

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft/meter"
	"github.com/Aj002Th/BlockchainEmulator/data/base"
	"github.com/Aj002Th/BlockchainEmulator/data/chain"
	"github.com/Aj002Th/BlockchainEmulator/misc"
	"github.com/Aj002Th/BlockchainEmulator/network"
	"github.com/Aj002Th/BlockchainEmulator/params"
	"github.com/Aj002Th/BlockchainEmulator/storage/blockStorage"
	"github.com/Aj002Th/BlockchainEmulator/storage/stateStorage"
)

type ConsensusNode struct {
	RunningNode *Node
	ShardID     uint64
	NodeID      uint64

	CurChain *chain.BlockChain
	db       blockStorage.BlockStorage
	sb       stateStorage.StateStorage

	pbftChainConfig *chain.Config
	ipNodeTable     map[uint64]string
	nodeNums        uint64
	maliciousNums   uint64
	view            uint64

	sequenceID        uint64
	stop              bool
	pStop             chan uint64
	requestPool       map[string]*Request
	cntPrepareConfirm map[string]map[*Node]bool
	cntCommitConfirm  map[string]map[*Node]bool
	isCommitBroadcast map[string]bool
	isReply           map[string]bool
	height2Digest     map[uint64]string

	sequenceLock sync.Mutex
	lock         sync.Mutex
	askForLock   sync.Mutex
	stopLock     sync.Mutex

	seqIDMap   map[uint64]uint64
	seqMapLock sync.Mutex

	pl *misc.PbftLog

	pbftImpl PbftImplInterface
}

func NewPbftNode(nodeID uint64, pcc *chain.Config) *ConsensusNode {
	self := new(ConsensusNode)
	self.ipNodeTable = params.IPmapNodeTable
	self.nodeNums = pcc.NodesNum
	self.ShardID = 0
	self.NodeID = nodeID
	self.pbftChainConfig = pcc

	var err error
	self.db = blockStorage.NewBoltStorage(uint(nodeID))
	self.sb = stateStorage.NewMemKVStore()
	self.CurChain, err = chain.NewBlockChain(pcc, self.sb, self.db)
	if err != nil {
		log.Panic("cannot new a blockchain")
	}

	self.RunningNode = &Node{
		NodeID:  nodeID,
		ShardID: 0,
		IPaddr:  self.ipNodeTable[nodeID],
	}

	self.stop = false
	self.sequenceID = self.CurChain.CurrentBlock.Header.Number + 1
	self.pStop = make(chan uint64)
	self.requestPool = make(map[string]*Request)
	self.cntPrepareConfirm = make(map[string]map[*Node]bool)
	self.cntCommitConfirm = make(map[string]map[*Node]bool)
	self.isCommitBroadcast = make(map[string]bool)
	self.isReply = make(map[string]bool)
	self.height2Digest = make(map[uint64]string)
	self.maliciousNums = (self.nodeNums - 1) / 3
	self.view = 0
	self.seqIDMap = make(map[uint64]uint64)
	self.pl = misc.NewPbftLog(0, nodeID)

	base.NodeLog = self.pl

	self.pbftImpl = &PbftImplSimpleImpl{
		node: self,
	}

	return self
}

func dispatchHelper[T any](content []byte, handler func(*T)) {
	message := new(T)
	err := json.Unmarshal(content, message)
	if err != nil {
		log.Panic(err)
	}
	handler(message)
}

func (self *ConsensusNode) dispatchMessage(msg []byte) {
	msgType, content := SplitMessage(msg)
	if len(content) > 2000 {
		self.pl.Printf("Received a %v: %v\n", msgType, string(content[:2000]))
	} else {
		self.pl.Printf("Received a %v: %v\n", msgType, string(content))
	}
	switch msgType {
	case CPrePrepare:
		dispatchHelper(content, self.handlePrePrepare)
	case CPrepare:
		dispatchHelper(content, self.handlePrepare)
	case CCommit:
		dispatchHelper(content, self.handleCommit)
	case CStop:
		self.setStopAndCleanUp()

	default:
		panic("Invalid Message Type.")
	}
}

func (self *ConsensusNode) Run() {
	meter.NodeSideStart()
	if self.NodeID == 0 {

		println("Node_id is 0 , so will do a propose")
		go self.doPropose()

		go func() {
			for {
				b := KeepAliveMsg{Msg: "alive"}
				m, err := json.Marshal(b)
				if err != nil {
					panic(err)
				}
				MergeAndSend(CKeepAlive, m, params.SupervisorEndpoint, self.pl)
				time.Sleep(time.Second * 5)
			}
		}()
	}
	self.serve()
}

func (self *ConsensusNode) serve() {
	ch := network.Tcp.Serve(self.RunningNode.IPaddr)
	for {
		clientRequest, ok := <-ch
		if !ok {
			self.pl.Println("From Pbft Node: the Tcp Channel is Closed")
			return
		}
		self.dispatchMessage(clientRequest)
	}
}

func (self *ConsensusNode) setStopAndCleanUp() {
	self.pl.Println("handling stop message")
	self.stopLock.Lock()
	self.stop = true
	self.stopLock.Unlock()
	if self.NodeID == self.view {
		self.pStop <- 1
	}

	self.pl.Println("Before GatherAndSend")

	GatherAndSend(int(self.NodeID), self.pl)

	network.Tcp.Close()
	self.CurChain.CloseBlockChain()
	self.pl.Println("handled stop message")
}

func GatherAndSend(nodeID int, pl *log.Logger) {

	b := BookingMsg{
		AvgCpuTime:    meter.AvgCpuPercent,
		DiskMetric:    meter.DiskMetric,
		TotalUpload:   meter.TotalUpload,
		TotalDownload: meter.TotalDownload,
		TotalTime:     int64(time.Since(meter.TimeBegin)),
		NodeId:        nodeID,
	}
	m, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	MergeAndSend(CBooking, m, params.SupervisorEndpoint, pl)
}

func (self *ConsensusNode) doPropose() {
	if self.view != self.NodeID {
		return
	}
	for {
		select {
		case <-self.pStop:
			self.pl.Printf("S%dN%d stop...\n", self.ShardID, self.NodeID)
			return
		default:
		}
		time.Sleep(time.Duration(int64(params.BlockInterval)) * time.Millisecond)

		self.sequenceLock.Lock()
		self.pl.Printf("S%dN%d get sequenceLock locked, now trying to propose...\n", self.ShardID, self.NodeID)

		_, r := self.pbftImpl.doPropose()

		digest := getDigest(r)
		self.requestPool[string(digest)] = r
		self.pl.Printf("S%dN%d put the request into the pool ...\n", self.ShardID, self.NodeID)

		ppmsg := PrePrepare{
			RequestMsg: r,
			Digest:     digest,
			SeqID:      self.sequenceID,
		}
		self.height2Digest[self.sequenceID] = string(digest)

		ppbyte, err := json.Marshal(ppmsg)
		if err != nil {
			log.Panic()
		}
		MergeAndBroadcast(CPrePrepare, ppbyte, self.RunningNode.IPaddr, self.getNeighborNodes(), self.pl)
	}
}

func MergeAndBroadcast(t MessageType, data []byte, from string, to []string, logger *log.Logger) {
	if len(data) < 2000 {
		logger.Printf("Broadcasting a %v: %v", t, string(data))
	} else {
		logger.Printf("Broadcasting a %v: %v", t, string(data[:2000]))
	}
	msgSend := MergeMessage(t, data)
	network.Tcp.Broadcast(from, to, msgSend)
}
