// bring up two nodes and connect them
package main

import (
	"fmt"
	"log"

	"github.com/anthoai97/devp2p-hand-on/utils"
	"github.com/ethereum/go-ethereum/p2p"
)

func main() {
	// we need private keys for both servers
	privkey_one := utils.GenerateKey()
	privkey_two := utils.GenerateKey()

	srv1 := &p2p.Server{Config: p2p.Config{
		PrivateKey:  privkey_one,
		MaxPeers:    1,
		NoDiscovery: true,
		Logger:      utils.Log.New("server", "number-1"),
	}}

	srv2 := &p2p.Server{Config: p2p.Config{
		PrivateKey:  privkey_two,
		MaxPeers:    1,
		NoDiscovery: true,
		NoDial:      true,
		ListenAddr:  fmt.Sprintf(":%d", 31234),
		Logger:      utils.Log.New("server", "number-2"),
	}}

	srv1.Start()
	defer srv1.Stop()
	srv2.Start()
	defer srv2.Stop()

	if !utils.SyncAddPeer(srv1, srv2.Self()) {
		log.Fatal("peer not connected")
	}

	utils.Log.Info("after add", "node one peers", srv1.Peers(), "node two peers", srv2.Peers())

	srv1.RemovePeer(srv2.Self())
	if srv1.PeerCount() > 0 {
		log.Fatal("removed peer still connected")
	}
}
