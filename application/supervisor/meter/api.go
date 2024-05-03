package meter

import (
	"fmt"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/measure"
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/metrics"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
	"github.com/Aj002Th/BlockchainEmulator/params"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Booking = pbft.BookingMsg

func SupSideStart() []*measure.MeasureModule {
	return StartCommitRelate()
}

func fmtBytes(v interface{}) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d bytes", v)
}

func GetResult(bookings *[]Booking) []metrics.Desc { // 每一个度量，作为一棵树，都是一个Desc。现在需要一系列Desc。
	var ds = make([]metrics.Desc, 0)

	// 统计TxCount和BlockNum和运行时间
	tx := metrics.NewDescBuilder("CPU时间", "交易计数，是指对交易的计数。")
	bc := metrics.NewDescBuilder("内存测量", "交易计数，是指对交易的计数。")
	dur := metrics.NewDescBuilder("墙上时钟时间", "")
	net := metrics.NewDescBuilder("网络", "")
	var sumC float64 = 0
	var sumBc uint64 = 0
	var sumDur time.Duration = 0
	var sumUp, sumDown int = 0, 0
	p := message.NewPrinter(language.English)
	for _, booking := range *bookings {
		nn := booking.NodeId
		c := booking.AvgCpuTime
		b := booking.DiskMetric
		t := time.Duration(booking.TotalTime)
		bStr := p.Sprintf("%d bytes", b)
		tx.AddElem(fmt.Sprintf("节点%v CPU负载", nn), "本节点的CPU负载", fmt.Sprintf("%.1f %%", c))
		bc.AddElem(fmt.Sprintf("节点%v 内存占用", nn), "本节点的内存占用", bStr)
		dur.AddElem(fmt.Sprintf("节点%v 墙上时钟时间", nn), "本节点的墙上时钟时间", fmt.Sprintf("%v min %v sec", int(t.Minutes()), int(t.Seconds())-int(t.Minutes())*60))
		net.AddElem(fmt.Sprintf("节点%v 上传流量", nn), "本节点的上传流量", fmtBytes(booking.TotalUpload))
		net.AddElem(fmt.Sprintf("节点%v 下载流量", nn), "本节点的下载流量", fmtBytes(booking.TotalDownload))
		sumC += c
		sumBc += b
		sumDur += t
		sumUp += booking.TotalUpload
		sumDown += booking.TotalDownload
	}
	tt := sumDur / time.Duration(params.NodeNum)

	// 系统指标
	bStr := p.Sprintf("%d bytes", sumBc/uint64(params.NodeNum))
	tx.AddElem("平均CPU负载", "平均的CPU负载", fmt.Sprintf("%.1f %%", sumC/float64(params.NodeNum)))
	bc.AddElem("平均内存用量", "平均的内存用量", bStr)
	dur.AddElem("平均运行时间（墙上时钟）", "平均运行时间，基准是Wall Time", fmt.Sprintf("%v min %v sec", int(tt.Minutes()), int(tt.Seconds())-int(tt.Minutes())*60))
	net.AddElem("总计上传流量", "总计上传流量", fmtBytes(sumUp))
	net.AddElem("总计下载流量", "总计下载流量", fmtBytes(sumDown))
	net.AddElem("平均上传流量", "平均上传流量", fmtBytes(sumUp/params.NodeNum))
	net.AddElem("平均下载流量", "平均下载流量", fmtBytes(sumDown/params.NodeNum))
	net.AddElem("总计流量", "总计流量", fmtBytes(sumUp+sumDown))
	net.AddElem("平均流量", "平均流量", fmtBytes((sumUp+sumDown)/params.NodeNum))
	ds = append(ds, tx.GetDesc())
	ds = append(ds, bc.GetDesc())
	ds = append(ds, dur.GetDesc())
	ds = append(ds, net.GetDesc())

	// 区块链指标
	ds = append(ds, mTCL.GetDesc())
	ds = append(ds, mAvgTPS.GetDesc())
	ds = append(ds, mPCL.GetDesc())
	ds = append(ds, mBlockNumCount.GetDesc())
	ds = append(ds, mTxNumCount.GetDesc())

	return ds
}
