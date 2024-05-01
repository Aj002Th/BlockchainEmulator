package meter

import "github.com/Aj002Th/BlockchainEmulator/network"

var TotalUpload int
var TotalDownload int

func StartNet() {
	network.Tcp.GetOnUpload().Connect(func(data int) {
		TotalUpload += data
	})
	network.Tcp.GetOnDownload().Connect(func(data int) {
		TotalDownload += data
	})
}
