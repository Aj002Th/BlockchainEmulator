package measure

import "github.com/Aj002Th/BlockchainEmulator/consensus/pbft"

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
	default:
		return nil
	}
}

type MeasureModule interface {
	UpdateMeasureRecord(*pbft.BlockInfoMsg)
	HandleExtraMessage([]byte)
	OutputMetricName() string
	OutputRecord() ([]float64, float64)
}
