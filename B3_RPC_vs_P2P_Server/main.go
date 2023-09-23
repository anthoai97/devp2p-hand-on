package main

import (
	"os"
	"runtime"

	"github.com/Microsoft/go-winio"
	"github.com/anthoai97/devp2p-hand-on/utils"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
)

// Get nodeinfo through RPC Server
// This example explain basic relationship between RPC Server and P2P Server
// in Ethereum blockchain
func main() {
	privkey_one := utils.GenerateKey()

	srv1 := p2p.Server{Config: p2p.Config{
		Name:        "server-number-1",
		PrivateKey:  privkey_one,
		MaxPeers:    1,
		NoDiscovery: true,
		Logger:      utils.Log.New("server", "number-1"),
	}}

	err := srv1.Start()
	if err != nil {
		utils.Log.Crit("Start p2p.Server #1 failed", "err", err)
	}
	defer srv1.Stop()

	// set up the RPC server
	server := rpc.NewServer()
	err = server.RegisterName("rpc", &srv1)
	if err != nil {
		utils.Log.Crit("Register API method(s) failed", "err", err)
	}

	// Listen on a random endpoint.
	endpoint := "go-ethereum-test-ipc"
	if runtime.GOOS == "windows" {
		endpoint = `\\.\pipe\` + endpoint
	} else {
		endpoint = os.TempDir() + "/" + endpoint
	}
	utils.Log.Info("endpoint", "data", endpoint)

	l, err := winio.ListenPipe(endpoint, nil)
	if err != nil {
		utils.Log.Crit("can't listen:", "err", err)
	}

	go server.ServeListener(l)
	defer server.Stop()

	client, err := rpc.Dial(endpoint)
	if err != nil {
		utils.Log.Crit("can't dial:", "err", err)
	}

	// Call get p2p nodeInfo
	var nodeInfo p2p.NodeInfo
	err = client.Call(&nodeInfo, "rpc_nodeInfo")
	if err != nil {
		utils.Log.Crit("RPC call fail", "err", err)
	}

	utils.Log.Info("server started", "enode", nodeInfo.Enode, "name", nodeInfo.Name, "ID", nodeInfo.ID, "IP", nodeInfo.IP)
	utils.Log.Info("Stop example process")
}
