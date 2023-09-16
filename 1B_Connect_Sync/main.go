// bring up two nodes and connect them
package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/anthoai97/devp2p-hand-on/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

// create a server
func newServer(privkey *ecdsa.PrivateKey, name string, version string, port int) *p2p.Server {

	// we need to explicitly allow at least one peer
	// otherwise the connection attempt will be refused
	cfg := p2p.Config{
		PrivateKey: privkey,
		Name:       name,
		MaxPeers:   2,
	}
	if port > 0 {
		cfg.ListenAddr = fmt.Sprintf(":%d", port)
	}
	srv := &p2p.Server{
		Config: cfg,
	}
	return srv
}

func main() {
	// we need private keys for both servers
	privkey_one, err := crypto.GenerateKey()
	if err != nil {
		utils.Log.Crit("Generate private key #1 failed", "err", err)
	}
	privkey_two, err := crypto.GenerateKey()
	if err != nil {
		utils.Log.Crit("Generate private key #2 failed", "err", err)
	}

	srv1 := &p2p.Server{Config: p2p.Config{
		PrivateKey:  privkey_one,
		MaxPeers:    1,
		NoDiscovery: true,
		Logger:      utils.Log.New("server", "1"),
	}}
	srv2 := &p2p.Server{Config: p2p.Config{
		PrivateKey:  privkey_two,
		MaxPeers:    1,
		NoDiscovery: true,
		NoDial:      true,
		ListenAddr:  "127.0.0.1:0",
		Logger:      utils.Log.New("server", "2"),
	}}
	srv1.Start()
	defer srv1.Stop()
	srv2.Start()
	defer srv2.Stop()

	s := strings.Split(srv2.ListenAddr, ":")
	if len(s) != 2 {
		log.Fatal("invalid ListenAddr")
	}
	// if port, err := strconv.Atoi(s[1]); err == nil {
	// 	srv2.localnode.Set(enr.TCP(uint16(port)))
	// }
	if !syncAddPeer(srv1, srv2.Self()) {
		log.Fatal("peer not connected")
	}

	utils.Log.Info("after add", "node one peers", srv1.Peers(), "node two peers", srv2.Peers())
	srv1.RemovePeer(srv2.Self())
	if srv1.PeerCount() > 0 {
		log.Fatal("removed peer still connected")
	}
}

func syncAddPeer(srv *p2p.Server, node *enode.Node) bool {
	var (
		ch      = make(chan *p2p.PeerEvent)
		sub     = srv.SubscribeEvents(ch)
		timeout = time.After(2 * time.Second)
	)
	defer sub.Unsubscribe()
	srv.AddPeer(node)
	for {
		select {
		case ev := <-ch:
			if ev.Type == p2p.PeerEventTypeAdd && ev.Peer == node.ID() {
				utils.Log.Info("Innn add", "node one peers", srv.Peers())
				return true
			}
		case <-timeout:
			return false
		}
	}
}
