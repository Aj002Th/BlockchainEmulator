package misc

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Aj002Th/BlockchainEmulator/logger"
	"github.com/Aj002Th/BlockchainEmulator/params"
)

type PbftLog = log.Logger

func NewPbftLog(sid, nid uint64) *PbftLog {
	prefix := fmt.Sprintf("S%dN%d: ", sid, nid)
	dirPath := params.LogWritePath + "/S" + strconv.Itoa(int(sid))
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	filePath := dirPath + "/N" + strconv.Itoa(int(nid)) + ".log"
	return logger.NewLogger(filePath, prefix)
}
