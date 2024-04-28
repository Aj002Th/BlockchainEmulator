package params

var (
	IPmap_nodeTable  = make(map[uint64]map[uint64]string)
	MeasureRelayMod  = []string{"TPS_Relay", "TCL_Relay", "CrossTxRate_Relay", "TxNumberCount_Relay"}
	Enabled_Measures = []string{"CPU", "Disk", "TPS", "TCL", "TxCount"}
	Meter_CPU        = true // 控制输出的度量
	Meter_Disk       = true
	Meter_TPS        = true
	Meter_TCL        = true
	Meter_TxCount    = true
)
