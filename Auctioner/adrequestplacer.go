package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

type adRequestPostBody struct {
	AdPlacementID string `json:"ad_placement_id"`
}

func placeAdRequest(bidder string, adPlacementID string, result chan<- Bid) {
	timeout := time.Duration(200 * time.Millisecond)
	client := http.Client{
		Timeout: timeout,
	}

	postBody := adRequestPostBody{AdPlacementID: adPlacementID}
	marshalledBody, _ := json.Marshal(postBody)
	payload := bytes.NewReader(marshalledBody)

	resp, err := client.Post("http://"+bidder+"/adrequest", "application/json", payload)
	if err != nil {
		errObj := err.(*url.Error)
		if errObj.Timeout() {
			fmt.Println(fmt.Sprintf("%s bidder was shorted", bidder))
		} else {
			fmt.Println(err)
		}
		return
	}
	// No bid case
	if resp.StatusCode == http.StatusNoContent {
		result <- Bid{}
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
