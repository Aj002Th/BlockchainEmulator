package params

// 系统的默认配置

var (
	BlockInterval      = 5000  // generate new block interval
	MaxBlockSizeGlobal = 2000  // the block contains the maximum number of transactions
	InjectSpeed        = 2000  // the transaction inject speed
	TotalDataSize      = 16000 // the total number of txs
	BatchSize          = 16000 // supervisor read a batch of txs then send them, it should be larger than inject speed
	NodeNum            = 3
	LogWritePath       = "./log"                  // log output path
	DataWritePath      = "./result"               // measurement data result output path
	RecordWritePath    = "./record"               // record output path
	SupervisorEndpoint = "127.0.0.1:18800"        //supervisor ip address
	FileInput          = `./BlockTransaction.csv` //the raw BlockTransaction data path
)
