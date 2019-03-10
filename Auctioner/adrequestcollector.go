package main

import (
	"time"
)

func collector(result <-chan Bid) []Bid {
	var bids []Bid
	done := make(chan bool)
	go runTimer(done)
	for {
		select {
		case bid := <-result:
			bids = append(bids, bid)
		case <-done:
			return bids
		}
	}
}

func runTimer(done chan<- bool) {
	timer := time.NewTimer(200 * time.Millisecond)
	<-timer.C
	// fmt.Println("timer expired")
	done <- true
}
