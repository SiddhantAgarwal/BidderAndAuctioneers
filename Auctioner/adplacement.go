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

	// Make channel buffered according to number of bidders
	result := make(chan Bid, len(Bidders))

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
	bids = filter(bids, func(bid Bid) bool {
		if bid.BidPrice == 0.0 && bid.AdID == "" {
			return false
		}
		return true
	})
	sort.Slice(bids, func(i, j int) bool {
		return bids[i].BidPrice > bids[j].BidPrice
	})
	if len(bids) > 0 {
		return &bids[0]
	}
	return nil
}

func filter(vs []Bid, f func(Bid) bool) []Bid {
	bids := make([]Bid, 0)
	for _, v := range vs {
		if f(v) {
			bids = append(bids, v)
		}
	}
	return bids
}
