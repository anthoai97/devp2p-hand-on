package main

import (
	"fmt"
	"os"

	"github.com/anthoai97/devp2p-hand-on/utils"
	"github.com/ethereum/go-ethereum/node"
)

var (
	p2pPort       = 30100
	ipcpath       = ".demo.ipc"
	datadirPrefix = ".data_"
)

func main() {
	// set up the service node
	cfg := &node.DefaultConfig
	cfg.P2P.ListenAddr = fmt.Sprintf(":%d", p2pPort)
	cfg.IPCPath = ipcpath
	cfg.DataDir = fmt.Sprintf("%s%d", datadirPrefix, p2pPort)

	// create the node instance with the config
	stack, err := node.New(cfg)
	if err != nil {
		utils.Log.Crit("ServiceNode create fail", "err", err)
	}

	// start the node
	err = stack.Start()
	if err != nil {
		utils.Log.Crit("ServiceNode start fail", "err", err)
	}
	defer os.RemoveAll(stack.DataDir())
	utils.Log.Info("ServiceNode start successful")

	// shut down
	err = stack.Close()
	if err != nil {
		utils.Log.Crit("Node stop fail", "err", err)
	}
	utils.Log.Info("ServiceNode closed successful")
}
