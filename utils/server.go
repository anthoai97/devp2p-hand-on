package utils

import (
	"crypto/ecdsa"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

func GenerateKey() *ecdsa.PrivateKey {
	privkey_one, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal("generate private key failed ", err)
	}

	return privkey_one
}

func SyncAddPeer(srv *p2p.Server, node *enode.Node) bool {
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
				return true
			}
		case <-timeout:
			return false
		}
	}
}
