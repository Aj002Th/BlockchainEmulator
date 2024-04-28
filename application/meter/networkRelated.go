package meter

import "github.com/Aj002Th/BlockchainEmulator/network"

var TotalUpload int
var TotalDownload int

func StartNet() {
	network.Tcp.OnUpload.Connect(func(data int) {
		TotalUpload += data
	})
	network.Tcp.OnDownload.Connect(func(data int) {
		TotalUpload += data
	})
}
