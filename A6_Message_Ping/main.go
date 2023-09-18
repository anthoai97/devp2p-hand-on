// bring up two nodes and connect them
package main

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/anthoai97/devp2p-hand-on/utils"
	"github.com/ethereum/go-ethereum/p2p"
)

// create a rotocol can send and receive "PING" "PONG" message each 10s
// the Run function is invoked upon connection
// it gets passed:
// * an instance of p2p.Peer, which represents the remote peer
// * an instance of p2p.MsgReadWriter, which is the io between the node and its peer
const (
	pingMsgCode = iota
	pongMsgCode
)

var (
	messageW = &sync.WaitGroup{}

	// Run implements the ping-pong protocol which sends ping messages to the peer
	// at 10s intervals, and responds to pings with pong messages.
	proto = p2p.Protocol{
		Name:    "ping-pong",
		Version: 1,
		Length:  1,
		Run: func(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
			log := utils.Log.New("peer.id", peer.ID())

			errC := make(chan error, 1)
			go func() {
				for range time.Tick(10 * time.Second) {
					log.Info("sending ping")
					if err := p2p.Send(rw, pingMsgCode, "PING"); err != nil {
						errC <- err
						return
					}
				}
			}()
			go func() {
				for {
					msg, err := rw.ReadMsg()
					if err != nil {
						errC <- err
						return
					}
					payload, err := io.ReadAll(msg.Payload)
					if err != nil {
						errC <- err
						return
					}
					log.Info("received message", "msg.code", msg.Code, "msg.payload", string(payload))
					if msg.Code == pingMsgCode {
						log.Info("sending pong")
						go p2p.Send(rw, pongMsgCode, "PONG")
					}
				}
			}()
			return <-errC
		},
	}
)

func main() {
	// we need private keys for both servers
	privkey_one := utils.GenerateKey()
	privkey_two := utils.GenerateKey()

	srv1 := &p2p.Server{Config: p2p.Config{
		PrivateKey:      privkey_one,
		MaxPeers:        1,
		NoDiscovery:     true,
		EnableMsgEvents: true,
		Protocols:       []p2p.Protocol{proto},
		Logger:          utils.Log.New("server", "number-1"),
	}}

	srv2 := &p2p.Server{Config: p2p.Config{
		PrivateKey:      privkey_two,
		MaxPeers:        1,
		NoDiscovery:     true,
		EnableMsgEvents: true,
		Protocols:       []p2p.Protocol{proto},
		NoDial:          true,
		ListenAddr:      fmt.Sprintf(":%d", 31234),
		Logger:          utils.Log.New("server", "number-2"),
	}}

	err := srv1.Start()
	if err != nil {
		utils.Log.Crit("Start p2p.Server #1 failed", "err", err)
	}

	err = srv2.Start()
	if err != nil {
		utils.Log.Crit("Start p2p.Server #2 failed", "err", err)
	}

	// set up the event subscriptions on both servers
	// the Err() on the Subscription object returns when subscription is closed

	// Setup Event subcription for server 1
	// eventOneC := make(chan *p2p.PeerEvent)
	// sub_one := srv1.SubscribeEvents(eventOneC)
	// messageW.Add(1)
	// // listen for event
	// go func() {
	// 	for {
	// 		peerEvent := <-eventOneC
	// 		if peerEvent.Type == p2p.PeerEventTypeAdd {
	// 			utils.Log.Debug("Received peer add notification on node #1", "peer", peerEvent.Peer)
	// 		} else if peerEvent.Type == p2p.PeerEventTypeMsgRecv {
	// 			utils.Log.Info("Received message nofification on node #1", "event", peerEvent)
	// 			messageW.Done()
	// 			return
	// 		}
	// 	}
	// }()

	// eventTwoC := make(chan *p2p.PeerEvent)
	// sub_two := srv2.SubscribeEvents(eventTwoC)
	// messageW.Add(1)
	// // listen for event
	// go func() {
	// 	for {
	// 		peerEvent := <-eventTwoC
	// 		if peerEvent.Type == p2p.PeerEventTypeAdd {
	// 			utils.Log.Debug("Received peer add notification on node #2", "peer", peerEvent.Peer)
	// 		} else if peerEvent.Type == p2p.PeerEventTypeMsgRecv {
	// 			utils.Log.Info("Received message nofification on node #2", "event", peerEvent)
	// 			messageW.Done()
	// 			return
	// 		}
	// 	}
	// }()

	if !utils.SyncAddPeer(srv1, srv2.Self()) {
		log.Fatal("peer not connected")
	}

	utils.Log.Info("after add", "node one peers", srv1.Peers(), "node two peers", srv2.Peers())
	//
	// wait for each respective message
	// messageW.Wait()

	// terminate subscription loops and unsubscribe
	// sub_one.Unsubscribe()
	// sub_two.Unsubscribe()

	// stop the servers
	defer srv1.Stop()
	defer srv2.Stop()

	select {}
}
