package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
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
		go placeAdRequest(bidder, result)
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

type Bid struct {
	Bidder   string
	BidPrice float64
	AdID     string
}

type AdObject struct {
	AdID     string  `json:"ad_id"`
	BidPrice float64 `json:"bid_price"`
}

func runTimer(done chan<- bool) {
	timer := time.NewTimer(200 * time.Millisecond)
	<-timer.C
	fmt.Println("timer expired")
	done <- true
}

func placeAdRequest(bidder string, result chan<- Bid) {
	timeout := time.Duration(200 * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}
	payload := strings.NewReader("{\n\t\"ad_placement_id\": \"qwerty\"\n}")
	resp, err := client.Post(bidder+"/adrequest", "application/json", payload)
	if err != nil {
		errObj := err.(*url.Error)
		if errObj.Timeout() {
			fmt.Println(fmt.Sprintf("%s bidder was shorted", bidder))
		}
		return
	}
	decoder := json.NewDecoder(resp.Body)
	var adObj AdObject
	err = decoder.Decode(&adObj)
	if err != nil {
		return
	}
	bid := Bid{
		Bidder:   bidder,
		BidPrice: adObj.BidPrice,
		AdID:     adObj.AdID,
	}
	result <- bid
}
