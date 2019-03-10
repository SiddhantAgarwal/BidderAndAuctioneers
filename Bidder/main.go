package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	fmt.Println("Let's place some bids")
	r := mux.NewRouter()

	registerRoutes(r)

	// port := ":" + os.Args[1]

	// Basic middlewares
	n := negroni.Classic()
	n.UseHandler(r)
	err := http.ListenAndServe(":80", n)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}
}

// IndexHandler : Index route handler
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	resMap := make(map[string]interface{})
	resp, _ := MakeJSONResponse("I am a bidder, sup dude!", resMap, true)
	SendJSONHttpResponse(w, resp, http.StatusOK)
}

func registerRoutes(router *mux.Router) {
	router.Methods("GET").Path("/").HandlerFunc(IndexHandler)
	router.Methods("POST").Path("/adrequest").HandlerFunc(AdRequestHandler)
}
