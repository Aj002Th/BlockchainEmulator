package boot

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"runtime/trace"

	"github.com/Aj002Th/BlockchainEmulator/params"
)

func getAbsPath() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

type App struct {
	args Args
}

func NewApp(a Args) App {
	return App{args: a}
}

func (self *App) Run() {
	var f *os.File
	var err error
	if self.args.isClient {
		f, err = os.Create("traceSup.out")
	} else {
		f, err = os.Create(fmt.Sprintf("traceNode%d.out", self.args.nodeID))
	}
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close trace file: %v", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	// 配置初始化, 最关键的内容是 IPmapNodeTable
	initConfig()

	if self.args.isClient {
		BuildSupervisor(self)
	} else {
		BuildNewPbftNode(uint64(self.args.nodeID), uint64(params.NodeNum))
	}
}
