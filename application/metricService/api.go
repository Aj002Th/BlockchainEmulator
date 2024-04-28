package metricservice

import (
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
	"github.com/Aj002Th/BlockchainEmulator/params"
	"github.com/Aj002Th/BlockchainEmulator/signal"
)

func StartByParams() {
	sig := signal.FindSignalByName[pbft.BlockInfoMsg]("OnCommit")
	sig.Connect(func(data pbft.BlockInfoMsg) {
		GatherAndSend(data)
	})
}

type MyInfo struct {
	cpuUsage float64
}

func GatherAndSend(bim pbft.BlockInfoMsg) {

}

type Service interface {
	Start()
	Gather(m *MyInfo)
}

func GetServiceByName(name string) Service {
	switch name {
	case "SVG":
		return &Cpu{}
	case "Disk":
		return &Mem{}
	default:
		panic("")
	}
}

func StartAll() {
	for _, v := range params.Enabled_Measures {
		s := GetServiceByName(v)
		s.Start()
	}
}
