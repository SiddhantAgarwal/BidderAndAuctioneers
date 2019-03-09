package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/rs/xid"
)

type adRequestPostBody struct {
	AdPlacementID string `json:"ad_placement_id"`
}

// AdObject : struct to wrap a AdObject
type AdObject struct {
	AdID     string  `json:"ad_id"`
	BidPrice float64 `json:"bid_price"`
}

func AdRequestHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t adRequestPostBody
	err := decoder.Decode(&t)
	respMap := make(map[string]interface{})
	if err != nil {
		respMap["error"] = "corrupt post body"
		resp, _ := MakeJSONResponse("corrupt post body", respMap, false)
		SendJSONHttpResponse(w, resp, http.StatusBadRequest)
		return
	}
	dice := rand.Intn(100)
	random := rand.NewSource(time.Now().UnixNano())
	if dice > 50 {
		resp, _ := json.Marshal(AdObject{
			AdID:     generateUniqueIDforAd(),
			BidPrice: rand.New(random).Float64() * 100,
		})
		SendJSONHttpResponse(w, resp, http.StatusOK)
	} else {
		SendJSONHttpResponse(w, nil, http.StatusNoContent)
	}
}

func generateUniqueIDforAd() string {
	x := xid.New()
	return x.String()
}
