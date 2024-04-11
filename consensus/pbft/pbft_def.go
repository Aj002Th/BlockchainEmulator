// The pbft consensus process

package pbft

import (
	"bufio"
	"io"
	"log"
	"net"
	"sync"

	"github.com/Aj002Th/BlockchainEmulator/data/chain"
	"github.com/Aj002Th/BlockchainEmulator/misc"
	"github.com/Aj002Th/BlockchainEmulator/network"
	"github.com/Aj002Th/BlockchainEmulator/params"
	"github.com/Aj002Th/BlockchainEmulator/storage/blockStorage"
	"github.com/Aj002Th/BlockchainEmulator/storage/stateStorage"
)

type PbftConsensusNode struct {
	// the local config about pbft
	RunningNode *Node  // the node information
	ShardID     uint64 // denote the ID of the shard (or pbft), only one pbft consensus in a shard
	NodeID      uint64 // denote the ID of the node in the pbft (shard)

	// the data structure for blockchain
	CurChain *chain.BlockChain // all node in the shard maintain the same blockchain
	db       blockStorage.BlockStorage
	sb       stateStorage.StateStorage

	// the global config about pbft
	pbftChainConfig *chain.Config                // the chain config in this pbft
	ip_nodeTable    map[uint64]map[uint64]string // denote the ip of the specific node
	node_nums       uint64                       // the number of nodes in this pfbt, denoted by N
	malicious_nums  uint64                       // f, 3f + 1 = N
	view            uint64                       // denote the view of this pbft, the main node can be inferred from this variant

	// the control message and message checking utils in pbft
	sequenceID        uint64                    // the message sequence id of the pbft
	stop              bool                      // send stop signal
	pStop             chan uint64               // channle for stopping consensus
	requestPool       map[string]*Request       // RequestHash to Request
	cntPrepareConfirm map[string]map[*Node]bool // count the prepare confirm message, [messageHash][Node]bool
	cntCommitConfirm  map[string]map[*Node]bool // count the commit confirm message, [messageHash][Node]bool
	isCommitBordcast  map[string]bool           // denote whether the commit is broadcast
	isReply           map[string]bool           // denote whether the message is reply
	height2Digest     map[uint64]string         // sequence (block height) -> request, fast read

	// locks about pbft
	sequenceLock sync.Mutex // the lock of sequence
	lock         sync.Mutex // lock the stage
	askForLock   sync.Mutex // lock for asking for a serise of requests
	stopLock     sync.Mutex // lock the stop varient

	// seqID of other Shards, to synchronize
	seqIDMap   map[uint64]uint64
	seqMapLock sync.Mutex

	// logger
	pl *misc.PbftLog
	// tcp control
	tcpln       net.Listener
	tcpPoolLock sync.Mutex

	// to handle the message in the pbft
	ihm ExtraOpInConsensus

	// to handle the message outside of pbft
	ohm OpInterShards
}

// generate a pbft consensus for a node
func NewPbftNode(shardID, nodeID uint64, pcc *chain.Config, messageHandleType string) *PbftConsensusNode {
	p := new(PbftConsensusNode)
	p.ip_nodeTable = params.IPmap_nodeTable
	p.node_nums = pcc.Nodes_perShard
	p.ShardID = shardID
	p.NodeID = nodeID
	p.pbftChainConfig = pcc
	// fp := "./record/ldb/s" + strconv.FormatUint(shardID, 10) + "/n" + strconv.FormatUint(nodeID, 10)
	var err error
	p.db = blockStorage.NewBoltStorage(uint(nodeID))
	p.sb = stateStorage.NewMemKVStore()
	if err != nil {
		log.Panic(err)
	}
	p.CurChain, err = chain.NewBlockChain(pcc, p.sb, p.db)
	if err != nil {
		log.Panic("cannot new a blockchain")
	}

	p.RunningNode = &Node{
		NodeID:  nodeID,
		ShardID: shardID,
		IPaddr:  p.ip_nodeTable[shardID][nodeID],
	}

	p.stop = false
	p.sequenceID = p.CurChain.CurrentBlock.Header.Number + 1
	p.pStop = make(chan uint64)
	p.requestPool = make(map[string]*Request)
	p.cntPrepareConfirm = make(map[string]map[*Node]bool)
	p.cntCommitConfirm = make(map[string]map[*Node]bool)
	p.isCommitBordcast = make(map[string]bool)
	p.isReply = make(map[string]bool)
	p.height2Digest = make(map[uint64]string)
	p.malicious_nums = (p.node_nums - 1) / 3
	p.view = 0

	p.seqIDMap = make(map[uint64]uint64)

	p.pl = misc.NewPbftLog(shardID, nodeID)

	// choose how to handle the messages in pbft or beyond pbft
	switch string(messageHandleType) {
	default:
		p.ihm = &RawRelayPbftExtraHandleMod{
			pbftNode: p,
		}
		p.ohm = &RawRelayOutsideModule{
			pbftNode: p,
		}
	}

	return p
}

// handle the raw message, send it to corresponded interfaces
func (p *PbftConsensusNode) handleMessage(msg []byte) {
	msgType, content := SplitMessage(msg)
	switch msgType {
	// pbft inside message type
	case CPrePrepare:
		p.handlePrePrepare(content)
	case CPrepare:
		p.handlePrepare(content)
	case CCommit:
		p.handleCommit(content)
	case CStop:
		p.WaitToStop()

	// handle the message from outside
	default:
		p.ohm.HandleMessageOutsidePBFT(msgType, content)
	}
}

func (p *PbftConsensusNode) handleClientRequest(con net.Conn) {
	defer con.Close()
	clientReader := bufio.NewReader(con)
	for {
		clientRequest, err := clientReader.ReadBytes('\n')
		if p.getStopSignal() {
			return
		}
		switch err {
		case nil:
			p.tcpPoolLock.Lock()
			p.handleMessage(clientRequest)
			p.tcpPoolLock.Unlock()
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}
	}
}

func (p *PbftConsensusNode) TcpListen() {
	ln, err := net.Listen("tcp", p.RunningNode.IPaddr)
	p.tcpln = ln
	if err != nil {
		log.Panic(err)
	}
	for {
		conn, err := p.tcpln.Accept()
		if err != nil {
			return
		}
		go p.handleClientRequest(conn)
	}
}

// when received stop
func (p *PbftConsensusNode) WaitToStop() {
	p.pl.Println("handling stop message")
	p.stopLock.Lock()
	p.stop = true
	p.stopLock.Unlock()
	if p.NodeID == p.view {
		p.pStop <- 1
	}
	network.Tcp.Close()
	p.tcpln.Close()
	p.closePbft()
	p.pl.Println("handled stop message")
}

func (p *PbftConsensusNode) getStopSignal() bool {
	p.stopLock.Lock()
	defer p.stopLock.Unlock()
	return p.stop
}

// close the pbft
func (p *PbftConsensusNode) closePbft() {
	p.CurChain.CloseBlockChain()
}
