package boot

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"runtime/trace"
	"time"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor"
	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/webapi"
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

	if self.args.isClient {
		if self.args.frontend {
			webapi.G_Proxy = webapi.NewGoodApiProxy()
			webapi.RunApiServer()
			webapi.RunFrontendServer()
			go exec.Command("start", "http://localhost:3000") // 把浏览器拉起来
		} else {
			webapi.G_Proxy = webapi.DumbProxy{}
		}

		webapi.G_Proxy.Enqueue(webapi.Hello)

		sup := supervisor.NewSupervisor()
		time.Sleep(10000 * time.Millisecond) // TODO: 去掉丑陋的Sleep
		sup.Run()
	} else {
		BuildNewPbftNode(uint64(self.args.nodeID), uint64(self.args.nodeNum))
	}
}
