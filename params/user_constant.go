package params

var (
	BlockInterval      = 5000  // generate new block interval
	MaxBlockSizeGlobal = 2000  // the block contains the maximum number of transactions
	InjectSpeed        = 2000  // the transaction inject speed
	TotalDataSize      = 16000 // the total number of txs
	BatchSize          = 16000 // supervisor read a batch of txs then send them, it should be larger than inject speed
	BrokerNum          = 10
	NodeNum            = 3
	DataWritePath      = "./result/"                               // measurement data result output path
	LogWritePath       = "./log"                                   // log output path
	RecordWritePath    = "./record"                                // record output path
	SupervisorEndpoint = "127.0.0.1:18800"                         //supervisor ip address
	FileInput          = `./2000000to2999999_BlockTransaction.csv` //the raw BlockTransaction data path
)
