package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Bid struct {
	Bidder   string  `json:"bidder"`
	BidPrice float64 `json:"bid_price"`
	AdID     string  `json:"ad_id"`
}

type AdObject struct {
	AdID     string  `json:"ad_id"`
	BidPrice float64 `json:"bid_price"`
}

func placeAdRequest(bidder string, adPlacementID string, result chan<- Bid) {
	timeout := time.Duration(200 * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}
	payload := strings.NewReader(fmt.Sprintf("{\n\t\"ad_placement_id\": \"%s\"\n}", adPlacementID))
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
