package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

type adPlacementPostBody struct {
	AdPlacementID string `json:"ad_placement_id"`
}

var Bidders []string = []string{"http://localhost:8081", "http://localhost:8082"}

func AdPlacementHandler(w http.ResponseWriter, r *http.Request) {
	respMap := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	var t adPlacementPostBody
	err := decoder.Decode(&t)
	if err != nil {
		respMap["error"] = "corrupt post body"
		resp, _ := MakeJSONResponse("corrupt post body", respMap, false)
		SendJSONHttpResponse(w, resp, http.StatusBadRequest)
		return
	}
	result := make(chan Bid)
	timeDone := make(chan bool)
	for _, bidder := range Bidders {
		go placeAdRequest(bidder, t.AdPlacementID, result)
	}
	go runTimer(timeDone)
	var bids []Bid
	for {
		select {
		case bid := <-result:
			bids = append(bids, bid)
		case <-timeDone:
			sort.Slice(bids, func(i, j int) bool {
				return bids[i].BidPrice > bids[j].BidPrice
			})
			if len(bids) > 0 {
				respMap["winner"] = bids[0]
			} else {
				respMap["winner"] = nil
			}
			resp, _ := MakeJSONResponse("Done", respMap, true)
			SendJSONHttpResponse(w, resp, http.StatusOK)
			return
		}
	}
}

func runTimer(done chan<- bool) {
	timer := time.NewTimer(200 * time.Millisecond)
	<-timer.C
	fmt.Println("timer expired")
	done <- true
}
