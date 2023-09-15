package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

// Create a p2p server
func newServer(name string, port int) (*p2p.Server, error) {
	pKey, err := crypto.GenerateKey()
	if err != nil {
		log.Printf("Generate private key failed with err: %v", err)
		return nil, err
	}
	cfg := p2p.Config{
		Name:       name,
		MaxPeers:   1,
		Protocols:  []p2p.Protocol{Proto},
		PrivateKey: pKey,
	}

	if port > 0 {
		cfg.ListenAddr = fmt.Sprintf(":%d", port)
	}

	srv := &p2p.Server{
		Config: cfg,
	}

	err = srv.Start()
	if err != nil {
		log.Printf("Start server failed with err: %v", err)
		return nil, err
	}

	return srv, nil
}

func connectPeer(svr *p2p.Server, enodeUrl string) error {
	// Parse enode url
	node, err := enode.Parse(enode.ValidSchemes, enodeUrl)

	if err != nil {
		log.Printf("Failed to parse enode url with err: %v", err)
		return err
	}

	svr.AddPeer(node)

	return nil
}
