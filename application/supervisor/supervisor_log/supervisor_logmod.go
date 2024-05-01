package supervisor_log

import (
	"io"
	"log"
	"os"

	"github.com/Aj002Th/BlockchainEmulator/params"
)

type SupervisorLog struct {
	Slog *log.Logger
}

func NewSupervisorLog() *SupervisorLog {
	writer1 := os.Stdout

	dirPath := params.LogWritePath
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	writer2, err := os.OpenFile(dirPath+"/Supervisor.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Panic(err)
	}
	pl := log.New(io.MultiWriter(writer1, writer2), "Supervisor: ", log.Lshortfile|log.Ldate|log.Ltime)
	return &SupervisorLog{
		Slog: pl,
	}
}

var Log1 = NewLogger1()

func NewLogger1() *log.Logger {
	writer1 := os.Stdout

	dirpath := params.LogWritePath
	err := os.MkdirAll(dirpath, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	writer2, err := os.OpenFile(dirpath+"/MySup.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Panic(err)
	}
	pl := log.New(io.MultiWriter(writer1, writer2), "My: ", log.Lshortfile|log.Ldate|log.Ltime)
	return pl
}
