package webapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Aj002Th/BlockchainEmulator/data/base"
)

// 用来解耦。GOlang貌似不能前置声明类+指针，所以用Interface做一下类型擦除Type Erasure。
type PostInterface interface {
	sendOneTx(tx base.Transaction)
}

type PostServer struct {
	pi *PostInterface
}

func NewPostServer(pi PostInterface) *PostServer {
	return &PostServer{
		pi: &pi,
	}
}

func (ps *PostServer) Run(pi *PostInterface) {
	log.SetFlags(0)
	http.HandleFunc("/api/status", ps.postTxHandler)
	http.ListenAndServe("0.0.0.0:17645", nil)
}

// 由于go的instance.method带类型擦除，所以可以把处理器写成一个方法。
// 毕竟函数是一等公民，Handler甚至可以是闭包。
func (ps *PostServer) postTxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	bs := make([]byte, 0)
	n, err := r.Body.Read(bs)
	if err != nil {
		panic("body read error in postApi")
	}
	if r.Method == "POST" {
		if r.Header.Get("Content-Type") == "application/json" {
			tx := new(base.Transaction)
			json.Unmarshal(bs[:n], &tx)
			(*ps.pi).sendOneTx(*tx)
		} else {
			panic("content-type error in postApi")
		}
	}

}
