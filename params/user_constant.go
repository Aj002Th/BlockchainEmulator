package params

// 系统的默认配置

var (
	NodeNum            int
	BlockInterval      int    // generate new block interval
	MaxBlockSizeGlobal int    // the block contains the maximum number of transactions
	InjectSpeed        int    // the transaction inject speed
	TotalDataSize      int    // the total number of txs
	BatchSize          int    // supervisor read a batch of txs then send them, it should be larger than inject speed
	LogWritePath       string // log output path
	DataWritePath      string // measurement data result output path
	RecordWritePath    string // record output path
	SupervisorEndpoint string //supervisor ip address
	FileInput          string //the raw BlockTransaction data path
)
