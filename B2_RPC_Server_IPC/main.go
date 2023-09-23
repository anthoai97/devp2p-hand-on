package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/Microsoft/go-winio"
	"github.com/anthoai97/devp2p-hand-on/utils"
	"github.com/ethereum/go-ethereum/rpc"
)

// create a rotocol can send and receive "PING" "PONG" message each 10s
// the Run function is invoked upon connection
// it gets passed:
// * an instance of p2p.Peer, which represents the remote peer
// * an instance of p2p.MsgReadWriter, which is the io between the node and its peer

func main() {
	server := rpc.NewServer()
	if err := server.RegisterName("ping", new(utils.PingService)); err != nil {
		log.Fatalf(err.Error())
	}
	defer server.Stop()

	// Listen on a random endpoint.
	endpoint := fmt.Sprintf("go-ethereum-test-ipc")
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

	// Create Client
	rpcClient, err := rpc.Dial(endpoint)
	if err != nil {
		utils.Log.Crit("can't dial:", "err", err)
	}

	var resp utils.PingResult
	err = rpcClient.Call(&resp, "ping_ping", "Ping")
	if err != nil {
		log.Fatal("read error:", err)
	}
	utils.Log.Info("Recevice RPC Response", "data", resp)
}
