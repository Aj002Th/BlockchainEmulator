package measure

import (
	"encoding/json"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/metrics"
	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
)

func GetByName(name string) MeasureModule {
	switch name {
	case "TPS_Pbft":
		return NewTestModule_avgTPS_Pbft()
	case "TCL_Pbft":
		return (NewTestModule_TCL_Pbft())
	case "CrossTxRate_Pbft":
		return (NewTestCrossTxRate_Pbft())
	case "TxNumberCount_Pbft":
		return (NewTestTxNumCount_Pbft())
	case "BlockNumCount":
		return (NewBlockNumCount())
	default:
		panic("Wrong Measure Name")
		return nil
	}
}

type MeasureModule interface {
	UpdateMeasureRecord(*pbft.BlockInfoMsg)
	OutputMetricName() string
	OutputRecord() ([]float64, float64)
	GetDesc() metrics.Desc
}

func MarshalDesc(d metrics.Desc) []byte {
	a, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return a
}

func PrintDescJson(d metrics.Desc) string {
	a := MarshalDesc(d)
	return string(a)
}
