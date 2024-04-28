package metrics

import (
	"fmt"

	"github.com/Aj002Th/BlockchainEmulator/application/comm"
	"github.com/Aj002Th/BlockchainEmulator/params"
	"github.com/Aj002Th/BlockchainEmulator/signal"
)

var metrics []Metrics

func StartHandle() {
	// 先把模块注册上来
	for _, v := range params.Enabled_Measures {
		m := GetMetricByName(v)
		metrics = append(metrics, m)
	}
	// 然后开始监听。听到就告诉各个模块处理。
	comm.Listen()
	m := comm.Recv()
	for _, v := range metrics {
		v.OnBlockInfoMsg(m)
	}
}

func GetMetricByName(name string) Metrics {
	switch name {
	case "TPS":
		return BlockCountMetric()
	default:
		panic("")
	}
}

func DumpAll() {
	for _, v := range metrics {
		desc := v.GetResult()
		fmt.Printf("desc: %v\n", desc)
	}
}

var OnDump signal.Signal[string]
