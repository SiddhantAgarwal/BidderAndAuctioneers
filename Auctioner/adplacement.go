package main

import (
	"encoding/json"
	"net/http"
	"sort"
)

type adPlacementPostBody struct {
	AdPlacementID string `json:"ad_placement_id"`
}

// var Bidders = []string{"http://bidder1", "http://bidder2"}
var Bidders []string

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

	// Dispatch requests to all bidders
	for _, bidder := range Bidders {
		go placeAdRequest(bidder, t.AdPlacementID, result)
	}

	bids := collector(result)
	respMap["winner"] = processWinner(bids)
	resp, _ := MakeJSONResponse("Winner", respMap, true)
	SendJSONHttpResponse(w, resp, http.StatusOK)
}

func processWinner(bids []Bid) *Bid {
	sort.Slice(bids, func(i, j int) bool {
		return bids[i].BidPrice > bids[j].BidPrice
	})
	if len(bids) > 0 {
		return &bids[0]
	}
	return nil
}
