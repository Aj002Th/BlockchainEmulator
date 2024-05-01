package committee

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/signal"
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/supervisor_log"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
	"github.com/Aj002Th/BlockchainEmulator/data/base"
	"github.com/Aj002Th/BlockchainEmulator/misc"
	"github.com/Aj002Th/BlockchainEmulator/params"
)

type RelayCommitteeModule struct {
	csvPath      string
	dataTotalNum int
	nowDataNum   int
	batchDataNum int
	IpNodeTable  map[uint64]map[uint64]string
	sl           *supervisor_log.SupervisorLog
	Ss           *signal.StopSignal // to control the stop message sending
}

var log1 = supervisor_log.Log1

func NewRelayCommitteeModule(Ip_nodeTable map[uint64]map[uint64]string, Ss *signal.StopSignal, slog *supervisor_log.SupervisorLog, csvFilePath string, dataNum, batchNum int) *RelayCommitteeModule {
	return &RelayCommitteeModule{
		csvPath:      csvFilePath,
		dataTotalNum: dataNum,
		batchDataNum: batchNum,
		nowDataNum:   0,
		IpNodeTable:  Ip_nodeTable,
		Ss:           Ss,
		sl:           slog,
	}
}

// transform, data to transaction
// check whether it is a legal txs message. if so, read txs and put it into the txList
func data2tx(data []string, nonce uint64) (*base.Transaction, bool) {
	if data[6] == "0" && data[7] == "0" && len(data[3]) > 16 && len(data[4]) > 16 && data[3] != data[4] {
		val, ok := new(big.Int).SetString(data[8], 10)
		if !ok {
			log.Panic("new int failed\n")
		}
		tx := base.NewTransaction(data[3][2:], data[4][2:], val, nonce)
		return tx, true
	}
	return &base.Transaction{}, false
}

// 把这一批交易发送给pbft主节点。
func (rthm *RelayCommitteeModule) txSending(txlist []*base.Transaction) {
	// the txs will be sent
	sendToShard := make(map[uint64][]*base.Transaction)

	// 把每个tx按顺序推进去
	for idx := 0; idx <= len(txlist); idx++ {
		// 到达InjectSpeed的时候就发一个InjectTxs给Node 0
		// 同时清空队列列表。
		if idx > 0 && (idx%params.InjectSpeed == 0 || idx == len(txlist)) {
			// send to shard
			it := pbft.InjectTxs{
				Txs:       sendToShard[0],
				ToShardID: 0,
			}
			itByte, err := json.Marshal(it)
			if err != nil {
				log.Panic(err)
			}
			pbft.MergeAndSend(pbft.CInject, itByte, rthm.IpNodeTable[0][0], log1)

			sendToShard = make(map[uint64][]*base.Transaction)
			time.Sleep(time.Second)
		}
		if idx == len(txlist) {
			break
		}
		tx := txlist[idx]
		senderSID := uint64(misc.Addr2Shard(tx.Sender))
		sendToShard[senderSID] = append(sendToShard[senderSID], tx)
	}
}

// MsgSendingControl
// Sup开启的时候同步地调一次。
// 把tx读出来然后用txSending发出去。
// read transactions, the Number of the transactions is - batchDataNum
func (rthm *RelayCommitteeModule) MsgSendingControl() {
	log1.Println("in MsgSendingControl")

	txFile, err := os.Open(rthm.csvPath)
	if err != nil {
		log.Panic(err)
	}
	defer txFile.Close()
	reader := csv.NewReader(txFile)
	txList := make([]*base.Transaction, 0) // save the txs in this epoch (round)

	// 将csv里每batchNum行转换成tx的列表，并且同步地发送出去。
	for {
		data, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}

		// 从csv row 转换到 base.Transaction
		if tx, ok := data2tx(data, uint64(rthm.nowDataNum)); ok {
			// 附加到txList
			txList = append(txList, tx)
			// 计数
			rthm.nowDataNum++
		}

		// re-shard condition, enough edges
		if len(txList) == rthm.batchDataNum || rthm.nowDataNum == rthm.dataTotalNum {
			// 当 txList 的长度满足batchDataNum，或者到达顶峰，那就要发送。
			rthm.txSending(txList)
			// reset the variants about tx sending
			txList = make([]*base.Transaction, 0) // 之后txList又恢复，从头再来。
			rthm.Ss.StopGap_Reset()               // StopGap也要Reset。
		}

		if rthm.nowDataNum == rthm.dataTotalNum {
			break
		}
	}
}

// HandleBlockInfo
// Sup会在HandleBlockInfos里调
// no operation here
func (rthm *RelayCommitteeModule) HandleBlockInfo(b *pbft.BlockInfoMsg) {
	// log1.Println("module HandleBlockInfo")
	// NOTHING TO DO
}
