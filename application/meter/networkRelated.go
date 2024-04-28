package meter

import "github.com/Aj002Th/BlockchainEmulator/network"

func onUpload(n int) {

}

func onDownload(n int) {

}

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
