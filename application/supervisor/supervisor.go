// Supervisor is an abstract role in this simulator that may read txs, generate partition infos,
// and handle history data.

package supervisor

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/committee"
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/measure"
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/signal"
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/supervisor_log"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
	"github.com/Aj002Th/BlockchainEmulator/data/chain"
	"github.com/Aj002Th/BlockchainEmulator/network"
	"github.com/Aj002Th/BlockchainEmulator/params"
	sig "github.com/Aj002Th/BlockchainEmulator/signal"
)

type Supervisor struct {
	// basic infos
	IPaddr       string // ip address of this Supervisor
	ChainConfig  *chain.Config
	Ip_nodeTable map[uint64]map[uint64]string

	// tcp control
	listenStop bool
	tcpLn      net.Listener
	tcpLock    sync.Mutex
	// logger module
	sl *supervisor_log.SupervisorLog

	// control components
	Ss *signal.StopSignal // to control the stop message sending

	// supervisor and committee components
	comMod committee.CommitteeModule

	// measure components
	testMeasureMods []measure.MeasureModule

	// 内部信号
	blockPostedSignal sig.Signal[pbft.BlockInfoMsg]

	// diy, add more structures or classes here ...
}

func (d *Supervisor) NewSupervisor(ip string, pcc *chain.Config, committeeMethod string, measureModNames ...string) {

	d.IPaddr = ip
	d.ChainConfig = pcc
	d.Ip_nodeTable = params.IPmap_nodeTable

	d.sl = supervisor_log.NewSupervisorLog()

	d.Ss = signal.NewStopSignal(2 * int(params.ShardNum))

	d.blockPostedSignal = sig.NewAsyncSignalImpl[pbft.BlockInfoMsg]("xx")

	switch committeeMethod {
	default:
		d.comMod = committee.NewRelayCommitteeModule(d.Ip_nodeTable, d.Ss, d.sl, params.FileInput, params.TotalDataSize, params.BatchSize)
	}

	d.testMeasureMods = make([]measure.MeasureModule, 0)

	for _, mModName := range measureModNames {
		m := measure.GetByName(mModName)
		d.testMeasureMods = append(d.testMeasureMods, m)
		d.blockPostedSignal.Connect(func(data pbft.BlockInfoMsg) { m.UpdateMeasureRecord(&data) })
	}

}

// Supervisor received the block information from the leaders, and handle these
// message to measure the performances.
func (d *Supervisor) handleBlockInfos(content []byte) {
	bim := new(pbft.BlockInfoMsg)
	err := json.Unmarshal(content, bim)
	if err != nil {
		log.Panic()
	}
	// StopSignal check
	if bim.BlockBodyLength == 0 {
		d.Ss.StopGap_Inc()
	} else {
		d.Ss.StopGap_Reset()
	}

	d.comMod.HandleBlockInfo(bim)

	// measure update
	d.blockPostedSignal.Emit(*bim)
	// add codes here ...
}

// read transactions from dataFile. When the number of data is enough,
// the Supervisor will do re-partition and send partitionMSG and txs to leaders.
func (d *Supervisor) SupervisorTxHandling() {
	d.comMod.MsgSendingControl()

	// TxHandling is end
	for !d.Ss.GapEnough() { // wait all txs to be handled
		time.Sleep(time.Second)
	}
	// send stop message
	stopmsg := pbft.MergeMessage(pbft.CStop, []byte("this is a stop message~"))
	d.sl.Slog.Println("Supervisor: now sending cstop message to all nodes")
	for sid := uint64(0); sid < uint64(params.ShardNum); sid++ {
		for nid := uint64(0); nid < d.ChainConfig.Nodes_perShard; nid++ {
			network.Tcp.Send(stopmsg, d.Ip_nodeTable[sid][nid])
		}
	}
	d.sl.Slog.Println("Supervisor: now Closing")
	d.listenStop = true
	d.CloseSupervisor()
}

// handle pbft. only one message to be handled now
func (d *Supervisor) handleMessage(msg []byte) {
	msgType, content := pbft.SplitMessage(msg)
	switch msgType {
	case pbft.CBlockInfo:
		d.handleBlockInfos(content)
		// add codes for more functionality
	default:
		d.comMod.HandleOtherMessage(msg)
		for _, mm := range d.testMeasureMods {
			mm.HandleExtraMessage(msg)
		}
	}
}

func (d *Supervisor) handleClientRequest(con net.Conn) {
	defer con.Close()
	clientReader := bufio.NewReader(con)
	for {
		clientRequest, err := clientReader.ReadBytes('\n')
		switch err {
		case nil:
			d.tcpLock.Lock()
			d.handleMessage(clientRequest)
			d.tcpLock.Unlock()
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}
	}
}

func (d *Supervisor) TcpListen() { // 不停听并且起goroutine
	ln, err := net.Listen("tcp", d.IPaddr)
	if err != nil {
		log.Panic(err)
	}
	d.tcpLn = ln
	for {
		conn, err := d.tcpLn.Accept()
		if err != nil {
			return
		}
		go d.handleClientRequest(conn)
	}
}

// close Supervisor, and record the data in .csv file
func (d *Supervisor) CloseSupervisor() {
	d.sl.Slog.Println("Closing...")
	for _, measureMod := range d.testMeasureMods {
		d.sl.Slog.Println(measureMod.OutputMetricName())
		d.sl.Slog.Println(measureMod.OutputRecord())
		println()
	}

	d.sl.Slog.Println("Trying to input .csv")
	// write to .csv file
	dirpath := params.DataWrite_path + "supervisor_measureOutput/"
	err := os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	for _, measureMod := range d.testMeasureMods {
		targetPath := dirpath + measureMod.OutputMetricName() + ".csv"
		f, err := os.Open(targetPath)
		resultPerEpoch, totResult := measureMod.OutputRecord()
		resultStr := make([]string, 0)
		for _, result := range resultPerEpoch {
			resultStr = append(resultStr, strconv.FormatFloat(result, 'f', 8, 64))
		}
		resultStr = append(resultStr, strconv.FormatFloat(totResult, 'f', 8, 64))
		if err != nil && os.IsNotExist(err) {
			file, er := os.Create(targetPath)
			if er != nil {
				panic(er)
			}
			defer file.Close()

			w := csv.NewWriter(file)
			title := []string{measureMod.OutputMetricName()}
			w.Write(title)
			w.Flush()
			w.Write(resultStr)
			w.Flush()
		} else {
			file, err := os.OpenFile(targetPath, os.O_APPEND|os.O_RDWR, 0666)

			if err != nil {
				log.Panic(err)
			}
			defer file.Close()
			writer := csv.NewWriter(file)
			err = writer.Write(resultStr)
			if err != nil {
				log.Panic()
			}
			writer.Flush()
		}
		f.Close()
		d.sl.Slog.Println(measureMod.OutputRecord())
	}

	network.Tcp.Close()
	d.tcpLn.Close()
}
