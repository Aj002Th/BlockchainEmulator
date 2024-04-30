package params

var (
	Block_Interval      = 5000  // generate new block interval
	MaxBlockSize_global = 2000  // the block contains the maximum number of transactions
	InjectSpeed         = 2000  // the transaction inject speed
	TotalDataSize       = 16000 // the total number of txs
	BatchSize           = 16000 // supervisor read a batch of txs then send them, it should be larger than inject speed
	BrokerNum           = 10
	NodeNum             = 3
	DataWrite_path      = "./result/"                                                           // measurement data result output path
	LogWrite_path       = "./log"                                                               // log output path
	RecordWrite_path    = "./record"                                                            // record output path
	SupervisorEndpoint  = "127.0.0.1:18800"                                                     //supervisor ip address
	FileInput           = `C:\Users\Administrator\block-dataset\0to999999_BlockTransaction.csv` //the raw BlockTransaction data path
)
