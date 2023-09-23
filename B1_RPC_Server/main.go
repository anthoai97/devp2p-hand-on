package main

import (
	"log"
	"net"

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

	utils.Log.Info("Start")
	defer server.Stop()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		utils.Log.Crit("can't listen:", err)
	}
	defer listener.Close()

	go server.ServeListener(listener)

	var (
		requestPing = `{
			"jsonrpc":"2.0",
			"id":1,
			"method":"ping_ping",
			"params" : [
			]
		}`
	)

	// Create connection
	utils.Log.Info("Listenter Addr " + listener.Addr().String())
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		log.Fatal("can't dial:", err)
	}

	// Write request
	conn.Write([]byte(requestPing))
	conn.(*net.TCPConn).CloseWrite()

	// Now try to get the response.
	buf := make([]byte, 2000)
	n, err := conn.Read(buf)
	conn.Close()

	if err != nil {
		log.Fatal("read error:", err)
	}

	resp := string(buf[:n])
	utils.Log.Info("Recevice RPC Response", "data", resp)
}
