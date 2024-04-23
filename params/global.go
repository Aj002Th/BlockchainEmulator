package params

var (
	IPmap_nodeTable = make(map[uint64]map[uint64]string)
	MeasureRelayMod = []string{"TPS_Relay", "TCL_Relay", "CrossTxRate_Relay", "TxNumberCount_Relay"}
)
