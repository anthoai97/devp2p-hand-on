package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/Microsoft/go-winio"
	"github.com/anthoai97/devp2p-hand-on/utils"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
)

// This example use RPC Client to connect 2 peer together
func main() {
	privkey_one := utils.GenerateKey()
	privkey_two := utils.GenerateKey()

	svr1 := p2p.Server{Config: p2p.Config{
		Name:        "server-number-1",
		PrivateKey:  privkey_one,
		MaxPeers:    1,
		NoDiscovery: true,
		Logger:      utils.Log.New("server", "number-1"),
	}}

	svr2 := p2p.Server{Config: p2p.Config{
		PrivateKey:  privkey_two,
		MaxPeers:    1,
		NoDiscovery: true,
		NoDial:      true,
		ListenAddr:  fmt.Sprintf(":%d", 31234),
		Logger:      utils.Log.New("server", "number-2"),
	}}

	err := svr1.Start()
	if err != nil {
		utils.Log.Crit("Start p2p.Server #1 failed", "err", err)
	}
	defer svr1.Stop()

	err = svr2.Start()
	if err != nil {
		utils.Log.Crit("Start p2p.Server #2 failed", "err", err)
	}
	defer svr2.Stop()

	// set up the RPC server
	server := rpc.NewServer()
	err = server.RegisterName("rpc", &svr1)
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
	utils.Log.Info("server #1 info", "enode", nodeInfo.Enode, "name", nodeInfo.Name, "ID", nodeInfo.ID, "IP", nodeInfo.IP)

	// Join peer
	err = client.Call(nil, "rpc_addPeer", svr2.Self())
	if err != nil {
		utils.Log.Crit("RPC call fail", "err", err)
	}

	// Wait for peer connected
	ch := make(chan *p2p.PeerEvent)
	sub := svr1.SubscribeEvents(ch)
	timeout := time.After(2 * time.Second)
	defer sub.Unsubscribe()

	for {
		select {
		case ev := <-ch:
			if ev.Type == p2p.PeerEventTypeAdd && ev.Peer == svr2.Self().ID() {
				utils.Log.Info("after add", "node one peers", svr1.Peers(), "node two peers", svr2.Peers())
				return
			}
		case <-timeout:
			utils.Log.Error("time-out")
			return
		}
	}
}
