package meter

import (
	"fmt"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/metrics"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
	"github.com/Aj002Th/BlockchainEmulator/params"
	"github.com/Aj002Th/BlockchainEmulator/signal"
	"github.com/chebyrash/promise"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Void = struct{}
type Booking = pbft.Booking

func SupSideStart() {
	StartCommitRelate()
}

var SupOnGathered signal.Signal[metrics.Desc]

func GetNReturnAsync() *promise.Promise[[]Booking] {
	sig := signal.GetSignalByName[Booking]("OnBookingReach")
	cnt := 0
	bks := make([]Booking, 0)
	p := promise.New(func(resolve func([]Booking), reject func(error)) {
		sig.Connect(func(bk Booking) {
			cnt++
			bks = append(bks, bk)
			if cnt == params.NodeNum {
				resolve(bks)
			}
		})
	})
	return p
}

func GetResult(ws *[]Booking) []metrics.Desc { // 每一个度量，作为一棵树，都是一个Desc。现在需要一系列Desc。
	// pws := GetNReturnAsync()
	// ws, _ := pws.Await(context.Background()) // 暂时没err，不用管err
	var ds = make([]metrics.Desc, 0)

	// 统计TxCount和BlockNum和运行时间
	tx := metrics.NewDescBuilder("CPU时间", "交易计数，是指对交易的计数。")
	bc := metrics.NewDescBuilder("内存测量", "交易计数，是指对交易的计数。")
	dur := metrics.NewDescBuilder("时间", "")
	net := metrics.NewDescBuilder("网络", "")
	var sumC float64 = 0
	var sumBc uint64 = 0
	var sumDur uint64 = 0
	var sumUp, sumDown int = 0, 0
	p := message.NewPrinter(language.English)
	for _, w := range *ws {
		nn := w.NodeId
		c := w.AvgCpuTime
		b := w.DiskMetric
		t := w.TotalTime
		bStr := p.Sprintf("%d bytes", b)
		tx.AddElem(fmt.Sprintf("节点%v CPU时间", nn), "", c)
		bc.AddElem(fmt.Sprintf("节点%v 内存测量", nn), "", bStr)
		dur.AddElem(fmt.Sprintf("节点%v 时间", nn), "", t)
		net.AddElem(fmt.Sprintf("节点%v 上传流量(bytes)", nn), "", w.TotalUpload)
		net.AddElem(fmt.Sprintf("节点%v 下载流量(bytes)", nn), "", w.TotalDownload)
		sumC += c
		sumBc += b
		sumDur += t
		sumUp += w.TotalUpload
		sumDown += w.TotalDownload
	}
	bStr := p.Sprintf("%d bytes", (sumBc / uint64(params.NodeNum)))
	tx.AddElem("总计CPU时间", "", sumC/float64(params.NodeNum))
	bc.AddElem("平均内存用量", "", bStr)
	dur.AddElem("平均运行时间（墙上时钟）", "", sumDur/uint64(params.NodeNum))
	net.AddElem("总计上传流量", "", sumUp/params.NodeNum)
	net.AddElem("总计下载流量", "", sumDown/params.NodeNum)
	ds = append(ds, tx.GetDesc())
	ds = append(ds, bc.GetDesc())
	ds = append(ds, dur.GetDesc())
	ds = append(ds, net.GetDesc())

	// 那堆传统模块
	ds = append(ds, m1.GetDesc())
	ds = append(ds, m2.GetDesc())
	ds = append(ds, m3.GetDesc())
	ds = append(ds, m4.GetDesc())
	ds = append(ds, m5.GetDesc())
	ds = append(ds, m6.GetDesc())
	return ds
}
