// bring up two nodes and connect them
package main

import (
	"fmt"
	"log"

	"github.com/anthoai97/devp2p-hand-on/utils"
	"github.com/ethereum/go-ethereum/p2p"
)

var (
	quitC = make(chan bool)
)

func main() {
	// we need private keys for both servers
	privkey_one := utils.GenerateKey()
	privkey_two := utils.GenerateKey()

	srv1 := p2p.Server{Config: p2p.Config{
		PrivateKey:  privkey_one,
		MaxPeers:    1,
		NoDiscovery: true,
		Logger:      utils.Log.New("server", "number-1"),
	}}

	srv2 := p2p.Server{Config: p2p.Config{
		PrivateKey:  privkey_two,
		MaxPeers:    1,
		NoDiscovery: true,
		NoDial:      true,
		ListenAddr:  fmt.Sprintf(":%d", 31234),
		Logger:      utils.Log.New("server", "number-2"),
	}}

	err := srv1.Start()
	if err != nil {
		utils.Log.Crit("Start p2p.Server #1 failed", "err", err)
	}

	err = srv2.Start()
	if err != nil {
		utils.Log.Crit("Start p2p.Server #2 failed", "err", err)
	}

	// Setup Event subcription for server 1
	eventC := make(chan *p2p.PeerEvent)
	sub_one := srv1.SubscribeEvents(eventC)

	// listen for event
	go func() {
		for {
			select {
			case peerEvent, ok := <-eventC:
				if !ok {
					return
				}
				utils.Log.Info("received peerevent", "type", peerEvent.Type, "peer", peerEvent.Peer)
				if peerEvent.Type == p2p.PeerEventTypeDrop {
					quitC <- true // Quit process when catch event drop
					return
				}
			case err := <-sub_one.Err():
				if err != nil {
					utils.Log.Error("Error in peer event subscription err", err)
					quitC <- true // Quit process when sub err
				}
				return
			}
		}
	}()

	if !utils.SyncAddPeer(&srv1, srv2.Self()) {
		log.Fatal("peer not connected")
	}

	utils.Log.Info("after add", "node one peers", srv1.Peers(), "node two peers", srv2.Peers())

	srv1.RemovePeer(srv2.Self())
	if srv1.PeerCount() > 0 {
		log.Fatal("removed peer still connected")
	}

	<-quitC
	sub_one.Unsubscribe()
	defer srv1.Stop()
	defer srv2.Stop()
}
