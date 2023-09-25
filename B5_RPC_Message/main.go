package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/Microsoft/go-winio"
	"github.com/anthoai97/devp2p-hand-on/utils"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	msgC    = make(chan string)
	msgWg   = &sync.WaitGroup{}
	protoWg = &sync.WaitGroup{}
)

type FooMsg struct {
	Content string
}

var protocol = p2p.Protocol{
	Name:    "foo",
	Version: 99,
	Length:  1,
	Run: func(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
		// only one of peer will fan out the message
		content, ok := <-msgC
		if ok {
			outMessage := &FooMsg{
				Content: content,
			}

			// sned the message
			err := p2p.Send(rw, 0, outMessage)
			if err != nil {
				return fmt.Errorf("send p2p message fail: %v", err)
			}
			utils.Log.Info("sending message", "peer", peer, "msg", outMessage)
		}

		// wait for the subcription to end

		// terminate protocal
		return nil
	},
}

type FooAPI struct {
	sent bool
}

func (api *FooAPI) SendMessage(content string) error {
	if api.sent {
		return fmt.Errorf("already sent")
	}
	msgC <- content
	close(msgC)
	api.sent = true
	return nil
}

func newRPCServer() (*rpc.Server, error) {
	server := rpc.NewServer()
	err := server.RegisterName("foo", &FooAPI{})
	if err != nil {
		return nil, fmt.Errorf("register API method(s) fail: %v", err)
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

	return server, nil
}

// create a protocol that can take care of message sending
// the Run function is invoked upon connection
// it gets passed:
// * an instance of p2p.Peer, which represents the remote peer
// * an instance of p2p.MsgReadWriter, which is the io between the node and its peer
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

	// set up the event subscriptions on both servers
	// the Err() on the Subscription object returns when subscription is closed
	evOneC := make(chan *p2p.PeerEvent)
	evTwoC := make(chan *p2p.PeerEvent)
	sub_one := srv1.SubscribeEvents(evOneC)
	defer sub_one.Unsubscribe()
	sub_two := srv2.SubscribeEvents(evTwoC)
	defer sub_two.Unsubscribe()

	// Lister for event chanel
	go func() {
		for {
			select {
			case peerEvent := <-evOneC:
				if peerEvent.Type == p2p.PeerEventTypeAdd {
					utils.Log.Info("Received peer add notification on node #1", "peer", peerEvent.Peer)
				} else if peerEvent.Type == p2p.PeerEventTypeMsgRecv {
					utils.Log.Info("Received message send notification on node #1", "event", peerEvent)
					msgWg.Done()
				}
			case <-sub_one.Err():
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case peerEvent := <-evTwoC:
				if peerEvent.Type == p2p.PeerEventTypeAdd {
					utils.Log.Info("Received peer add notification on node #2", "peer", peerEvent.Peer)
				} else if peerEvent.Type == p2p.PeerEventTypeMsgRecv {
					utils.Log.Info("Received message send notification on node #2", "event", peerEvent)
					msgWg.Done()
				}
			case <-sub_two.Err():
				return
			}
		}
	}()

	// Create rpc server
	// msgWg.Wait()
}
