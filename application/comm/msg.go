package comm

import "github.com/Aj002Th/BlockchainEmulator/consensus/pbft"

type UniversalMetricMsg struct {
	Bim      pbft.BlockInfoMsg
	CpuUsage int
	MemUsage int
}

type MM = UniversalMetricMsg
