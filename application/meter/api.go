package meter

import (
	"context"
	"fmt"

	"github.com/Aj002Th/BlockchainEmulator/application/comm"
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/metrics"
	"github.com/Aj002Th/BlockchainEmulator/params"
	"github.com/Aj002Th/BlockchainEmulator/signal"
	"github.com/chebyrash/promise"
)

type Wrapper = comm.Wrapper

func SupSideStart() {
	signal.GetSignalByName[Void]("OnSupStart").Connect(func(Void) {
		// Sup只需启动Commit相关，因为这些计算比较重。
		StartCommitRelate()
	})
}

func NodeSideStart() {
	// Node需要启动区块计数等等模块。
	signal.GetSignalByName[Void]("OnNodeStart").Connect(func(Void) {
		StartCnt()
		StartNet()
		StartTimeCnt()
		StartPs()
	})
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

func GetResult() []metrics.Desc { // 每一个度量，作为一棵树，都是一个Desc。现在需要一系列Desc。
	pws := GetNReturnAsync()
	ws, _ := pws.Await(context.Background()) // 暂时没err，不用管err
	var ds = make([]metrics.Desc, 0)

	// 统计TxCount和BlockNum和运行时间
	tx := metrics.NewDescBuilder("CPU时间", "交易计数，是指对交易的计数。")
	bc := metrics.NewDescBuilder("内存测量", "交易计数，是指对交易的计数。")
	dur := metrics.NewDescBuilder("时间", "")
	var sumC uint64 = 0
	var sumBc uint64 = 0
	var sumDur uint64 = 0
	for _, w := range *ws {
		nn := w.NodeId
		c := w.TxCount
		b := w.BlockCount
		t := w.TotalTime
		tx.AddElem(fmt.Sprintf("节点%v CPU事件", nn), "", c)
		bc.AddElem(fmt.Sprintf("节点%v 内存测量", nn), "", b)
		dur.AddElem(fmt.Sprintf("节点%v 时间", nn), "", t)
		sumC += c
		sumBc += b
		sumDur += t
	}
	tx.AddElem("平均计数", "", sumC/uint64(params.NodeNum))
	bc.AddElem("平均计数", "", sumBc/uint64(params.NodeNum))
	dur.AddElem("平均运行时间", "", sumDur/uint64(params.NodeNum))
	ds = append(ds, tx.GetDesc())
	ds = append(ds, bc.GetDesc())
	ds = append(ds, dur.GetDesc())

	// 那堆传统模块
	ds = append(ds, m1.GetDesc())
	ds = append(ds, m2.GetDesc())
	ds = append(ds, m3.GetDesc())
	ds = append(ds, m4.GetDesc())
	ds = append(ds, m5.GetDesc())
	ds = append(ds, m6.GetDesc())
	return ds
}
