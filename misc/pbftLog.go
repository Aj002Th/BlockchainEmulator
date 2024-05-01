package misc

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/Aj002Th/BlockchainEmulator/params"
)

// type PbftLog struct { // TODO: 把这个没用的wrapper改成type alias
// 	Plog *log.Logger
// }

type PbftLog = log.Logger

func NewPbftLog(sid, nid uint64) *PbftLog {
	pfx := fmt.Sprintf("S%dN%d: ", sid, nid)
	writer1 := os.Stdout

	dirpath := params.LogWritePath + "/S" + strconv.Itoa(int(sid))
	err := os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	writer2, err := os.OpenFile(dirpath+"/N"+strconv.Itoa(int(nid))+".log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Panic(err)
	}
	pl := log.New(io.MultiWriter(writer1, writer2), pfx, log.Lshortfile|log.Ldate|log.Ltime)
	fmt.Println()

	return pl
}
