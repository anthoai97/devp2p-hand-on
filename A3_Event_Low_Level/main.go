package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/event"
)

// This example shows how the return value of Send can be used for request/reply
// interaction between event consumers and producers.“
func main() {

	var feed event.Feed
	type ackedEvent struct {
		i   int
		ack chan<- struct{}
	}

	// Consumers wait for events on the feed and acknowledge processing
	done := make(chan struct{})
	defer close(done)
	for i := 0; i < 3; i++ {
		ch := make(chan ackedEvent, 1000)
		sub := feed.Subscribe(ch)
		go func() {
			defer sub.Unsubscribe()
			for {
				select {
				case ev := <-ch:
					fmt.Println(ev.i) // "process" the event
					ev.ack <- struct{}{}
				case <-done:
					return
				}
			}
		}()
	}

	for i := 0; i < 3; i++ {
		asksignal := make(chan struct{})
		n := feed.Send(ackedEvent{i, asksignal})
		for ack := 0; ack < n; ack++ {
			<-asksignal
		}
	}

}
