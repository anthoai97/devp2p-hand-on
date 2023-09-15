package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/p2p"
)

// p2p Protocol definition for sending and receiving a messag

var Proto = p2p.Protocol{
	Name:    "proto",
	Length:  1,
	Version: 1,
	Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
		message := "ping"

		// Send message to conected peer
		err := p2p.Send(rw, 0, message)
		if err != nil {
			return fmt.Errorf("send message fail: %v", err)
		}

		// Receive message from connected peer
		receive, err := rw.ReadMsg()
		if err != nil {
			return fmt.Errorf("read message fail: %v", err)
		}

		var myMessage string
		err = receive.Decode(message)
		if err != nil {
			return fmt.Errorf("decode message fail: %v", err)
		}

		fmt.Println("received message", string(myMessage))

		return nil
	},
}
