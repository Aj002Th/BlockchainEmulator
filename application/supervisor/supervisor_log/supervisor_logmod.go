package supervisor_log

import (
	"log"
	"os"

	"github.com/Aj002Th/BlockchainEmulator/logger"
	"github.com/Aj002Th/BlockchainEmulator/params"
)

type SupervisorLog struct {
	Slog *log.Logger
}

func NewSupervisorLog() *SupervisorLog {
	dirPath := params.LogWritePath
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	filePath := dirPath + "/Supervisor.log"
	return &SupervisorLog{
		Slog: logger.NewLogger(filePath, "Supervisor: "),
	}
}

var DebugLog *log.Logger

func NewLogger1() *log.Logger {
	dirPath := params.LogWritePath
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Println(dirPath)
		log.Panic(err)
	}
	filePath := dirPath + "/debug.log"
	return logger.NewLogger(filePath, "debug")
}
