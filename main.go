package main

import (
	"fmt"
	"os"

	"github.com/Aj002Th/BlockchainEmulator/boot"
)

func main() {
	app, err := boot.InitializeApp()
	if err != nil {
		fmt.Printf("failed to create app: %s\n", err)
		os.Exit(2)
	}
	app.Run()
}
