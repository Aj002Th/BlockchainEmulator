package boot

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
)

func getAbsPath() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

func batch_start(nodeNum, shardNum, modID int) {
	var absolute_path = getAbsPath()
	os.Chdir(absolute_path)
	// node = 1..n, shard = 0..n
	for i := 1; i < nodeNum; i++ {
		for j := 0; j < shardNum; j++ {
			cmdNormNodes := fmt.Sprintf("go run main.go -n %d -N %d -s %d -S %d -m %d \n\n", i, nodeNum, j, shardNum, modID)
			exec.Command(cmdNormNodes)
		}
	}
	// supervisor
	cmdSup := fmt.Sprintf("go run main.go -c -N %d -S %d -m %d \n\n", nodeNum, shardNum, modID)
	exec.Command(cmdSup)

	// node = 0, shrad = 0..n
	for j := 0; j < shardNum; j++ {
		cmdMainNodes := fmt.Sprintf("go run main.go -n 0 -N %d -s %d -S %d -m %d \n\n", nodeNum, j, shardNum, modID)
		exec.Command(cmdMainNodes)
	}
}

type App struct {
	args Args
}

func NewApp(a Args) App {
	return App{args: a}
}

func (self *App) Run() {
	if self.args.isGen {
		batch_start(self.args.nodeNum, self.args.shardNum, self.args.modID)
		return
	}

	if self.args.isClient {
		BuildSupervisor(uint64(self.args.nodeNum), uint64(self.args.shardNum), uint64(self.args.modID))
	} else {
		BuildNewPbftNode(uint64(self.args.nodeID), uint64(self.args.nodeNum), uint64(self.args.shardID), uint64(self.args.shardNum), uint64(self.args.modID))
	}
}
