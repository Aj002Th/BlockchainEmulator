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
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/committee"
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/measure"
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/meter"
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/metrics"
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/signal"
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/supervisor_log"
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/webapi"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
	"github.com/Aj002Th/BlockchainEmulator/network"
	"github.com/Aj002Th/BlockchainEmulator/params"
	sig "github.com/Aj002Th/BlockchainEmulator/signal"
)

type Supervisor struct {
	IpNodeTable       map[uint64]map[uint64]string  // basic infos
	tcpLn             net.Listener                  // tcp control
	tcpLock           sync.Mutex                    // listenStop bool
	sl                *supervisor_log.SupervisorLog // logger module
	Ss                *signal.StopSignal            // control components// to control the stop message sending
	cmt               committee.CommitteeModule     // supervisor and committee components
	testMeasureMods   []measure.MeasureModule       // measure components
	blockPostedSignal sig.Signal[pbft.BlockInfoMsg] // 内部信号

	txCompleteCount int
	pbftItems       []webapi.PbftItem
	OnNodeStart     sig.Signal[struct{}]

	cntBooking int
	bookings   []pbft.BookingMsg
	result     chan []metrics.Desc

	waitOnce        func()
	waitMasterReady chan struct{}
}

var log1 = supervisor_log.Log1

func NewSupervisor() *Supervisor {
	d := &Supervisor{}
	d.IpNodeTable = params.IPmapNodeTable
	d.sl = supervisor_log.NewSupervisorLog()
	d.Ss = signal.NewStopSignal(2 * int(1))
	d.blockPostedSignal = sig.NewAsyncSignalImpl[pbft.BlockInfoMsg]("xx")
	d.cmt = committee.NewPbftCommitteeModule(d.IpNodeTable, d.Ss, d.sl, params.FileInput, params.TotalDataSize, params.BatchSize)
	d.testMeasureMods = make([]measure.MeasureModule, 0)
	d.txCompleteCount = 0

	// 测量模块的附加。
	for _, mModName := range params.MeasurePbftMod {
		m := measure.GetByName(mModName)
		d.testMeasureMods = append(d.testMeasureMods, m)
		d.blockPostedSignal.Connect(func(data pbft.BlockInfoMsg) { m.UpdateMeasureRecord(&data) })
	}
	d.OnNodeStart = sig.GetSignalByName[struct{}]("OnNodeStart")

	d.result = make(chan []metrics.Desc)

	// 主节点 keep alive 相关处理
	d.waitMasterReady = make(chan struct{})
	d.waitOnce = sync.OnceFunc(func() {
		println("OnceFunc called")
		d.waitMasterReady <- struct{}{}
	})

	return d
}

// 根据Body长度决定是否Inc.
// 并且触发测量模块。
func (d *Supervisor) handleBlockInfoMsg(m *pbft.BlockInfoMsg) {
	log1.Println("in handleBlockInfos")

	log1.Println("StopSignal Check")

	// StopSignal check
	if m.BlockBodyLength == 0 {
		log1.Println("BodyLength == 0, Inc")
		d.Ss.StopGap_Inc()
	} else {
		log1.Println("BodyLength != 0, Reset")
		d.Ss.StopGap_Reset()
	}

	log1.Printf("received from shard %d in epoch %d.\n", m.SenderShardID, m.Epoch)

	d.txCompleteCount += len(m.ExcutedTxs)
	webapi.GlobalProxy.Enqueue(webapi.Computing(params.TotalDataSize, d.txCompleteCount))

	pbftItem := webapi.PbftItem{TxpoolSize: int(m.TxpoolSize), Tx: len(m.ExcutedTxs), Ctx: int(m.Pbft1TxNum)}
	d.pbftItems = append(d.pbftItems, pbftItem)
	// measure update
	d.blockPostedSignal.Emit(*m)
	// add codes here ...
}

func (d *Supervisor) handleBookingMsg(m *pbft.BookingMsg) {
	d.cntBooking++
	d.bookings = append(d.bookings, *m)
	d.sl.Slog.Printf("handleBookingMsg, cnt = %v\n", d.cntBooking)
	if d.cntBooking == params.NodeNum {
		result := meter.GetResult(&d.bookings)
		d.sl.Slog.Printf("handleBookingMsg now got result: %v\n", result)
		d.sl.Slog.Printf("handleBookingMsg now writing to channel\n")
		d.result <- result
	}
}

func (d *Supervisor) handleKeepAliveMsg(m *pbft.KeepAliveMsg) {
	d.waitOnce()
}

func (d *Supervisor) Wait() {
	<-d.waitMasterReady
	time.Sleep(time.Second)
}

func (d *Supervisor) Run() {
	webapi.GlobalProxy.Enqueue(webapi.Started)
	meter.SupSideStart()

	// 起一个听的循环
	go d.doAccept()

	d.Wait()

	// 无脑发送全部东西给主节点。
	d.cmt.MsgSendingControl()
	log1.Println("afterMsgSendingControl")

	// 发完之后开始准备在恰当时机发送Stop信息。
	for !d.Ss.GapEnough() { // wait all txs to be handled
		time.Sleep(time.Second)
	}
	// send stop message
	d.sl.Slog.Println("Supervisor: now sending cstop message to all nodes")

	for nid := uint64(0); nid < uint64(params.NodeNum); nid++ {
		log1.Printf("Sending a %v: %v\n", pbft.CStop, string([]byte("this is a stop message~")))
		pbft.MergeAndSend(pbft.CStop, []byte("this is a stop message~"), d.IpNodeTable[0][nid], log1)
	}

	d.sl.Slog.Println("Supervisor: now Closing. Now Generate Metrics Outputs.")

	d.generateOutputAndCleanUp()
}

// handle pbft. only one message to be handled now
func (d *Supervisor) dispatchMessage(msg []byte) {
	msgType, content := pbft.SplitMessage(msg)
	if len(content) > 2000 {
		log1.Printf("Received a %v: %v\n", msgType, string(content[:2000]))
	} else {
		log1.Printf("Received a %v: %v\n", msgType, string(content))
	}
	switch msgType {
	case pbft.CKeepAlive: // 用于确认主节点的启动情况
		m := new(pbft.KeepAliveMsg)
		err := json.Unmarshal(content, m)
		if err != nil {
			log.Panic()
		}
		d.handleKeepAliveMsg(m)

	case pbft.CBlockInfo: // 统计区块相关指标
		m := new(pbft.BlockInfoMsg)
		err := json.Unmarshal(content, m)
		if err != nil {
			log.Panic()
		}
		d.handleBlockInfoMsg(m)
		si := sig.GetSignalByName[*pbft.BlockInfoMsg]("OnBimReached")
		si.Emit(m)

	case pbft.CBooking: // 统计节点运行状态相关指标
		m := new(pbft.BookingMsg)
		err := json.Unmarshal(content, m)
		if err != nil {
			log.Panic()
		}
		d.handleBookingMsg(m)

	default:
		panic("Message Unsupported")
	}
}

func (d *Supervisor) startSession(con net.Conn) {
	defer con.Close()
	clientReader := bufio.NewReader(con)
	for {
		clientRequest, err := clientReader.ReadBytes('\n')
		switch err {
		case nil:
			d.tcpLock.Lock()
			d.dispatchMessage(clientRequest)
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

func (d *Supervisor) doAccept() { // 不停听并且起goroutine
	// ln, err := net.Listen("tcp", params.SupervisorEndpoint)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// d.tcpLn = ln
	// for {
	// 	conn, err := d.tcpLn.Accept()
	// 	if err != nil {
	// 		return
	// 	}
	// 	log1.Printf("Accepted the: %v. Now Start a session.\n", conn.RemoteAddr())
	// 	go d.startSession(conn)
	// }
	ch := network.Tcp.Serve(params.SupervisorEndpoint)
	for {
		clientRequest := <-ch
		log.Printf("Receiving %v", clientRequest)
		d.dispatchMessage(clientRequest)
	}
}

// close Supervisor, and record the data in .csv file
func (d *Supervisor) generateOutputAndCleanUp() {
	d.sl.Slog.Println("Closing...")
	for _, measureMod := range d.testMeasureMods {
		d.sl.Slog.Println(measureMod.OutputMetricName())
		d.sl.Slog.Println(measureMod.OutputRecord())
	}

	d.sl.Slog.Println("Before input .csv")
	// write to .csv file
	dirpath := path.Join(params.DataWritePath, "supervisor_measureOutput/")
	err := os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	// var measureItems []webapi.MeasureItem = make([]webapi.MeasureItem, 0)

	for _, measureMod := range d.testMeasureMods { // 遍历测试模组
		targetPath := path.Join(dirpath, measureMod.OutputMetricName()+".csv")
		f, err := os.Open(targetPath)
		resultPerEpoch, totResult := measureMod.OutputRecord()

		allResult := make([]float64, 0)
		allResult = append(allResult, totResult)
		allResult = append(allResult, resultPerEpoch...)

		// 附加到包里。
		// measureItems = append(measureItems, webapi.MeasureItem{Name: measureMod.OutputMetricName(), Desc: measure.PrintDescJson(measureMod.GetDesc()), Vals: allResult})

		// 对于文件则控制精度
		resultStr := make([]string, 0)
		for _, result := range resultPerEpoch {
			resultStr = append(resultStr, strconv.FormatFloat(result, 'f', 8, 64))
		}
		// 拼接
		resultStr = append(resultStr, strconv.FormatFloat(totResult, 'f', 8, 64))
		if err != nil && os.IsNotExist(err) { // 不存在则创建文件并写入
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
		} else { // 存在则直接写入文件
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

	// webapi.GlobalProxy.Enqueue(webapi.Completed(d.pbftItems, measureItems))

	d.sl.Slog.Println("Now waiting for Other Node Bookings and result")

	result := <-d.result
	d.sl.Slog.Println("result generated")

	webapi.GlobalProxy.Enqueue(webapi.Completed1(d.pbftItems, result))

	network.Tcp.Close()
	d.tcpLn.Close()
	webapi.GlobalProxy.Enqueue(webapi.Bye)
}
