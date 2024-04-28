package measure

import (
	"encoding/json"

	"github.com/Aj002Th/BlockchainEmulator/consensus/pbft"
)

func GetByName(name string) MeasureModule {
	switch name {
	case "TPS_Relay":
		return NewTestModule_avgTPS_Relay()
	case "TCL_Relay":
		return (NewTestModule_TCL_Relay())
	case "CrossTxRate_Relay":
		return (NewTestCrossTxRate_Relay())
	case "TxNumberCount_Relay":
		return (NewTestTxNumCount_Relay())
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
	GetDesc() Desc
}

func MarshalDesc(d Desc) []byte {
	a, b := json.Marshal(d)
	if b != nil {
		panic("PrintDesc")
	}
	return a
}

func PrintDescJson(d Desc) string {
	a := MarshalDesc(d)
	return string(a)
}
