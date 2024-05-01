package params

var (
	IPmapNodeTable  = make(map[uint64]map[uint64]string)
	MeasurePbftMod  = []string{"TPS_Pbft", "TCL_Pbft", "CrossTxRate_Pbft", "TxNumberCount_Pbft"}
	EnabledMeasures = []string{"CPU", "Disk", "TPS", "TCL", "TxCount"}
	Meter_CPU       = true // 控制输出的度量
	Meter_Disk      = true
	Meter_TPS       = true
	Meter_TCL       = true
	Meter_TxCount   = true
)
