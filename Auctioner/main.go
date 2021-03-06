package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	fmt.Println("I am a greedy Auctioner, let's roll")

	Bidders = strings.Split(os.Getenv("BIDDERS"), ",")

	// port := ":" + os.Args[1]

	// Basic middlewares
	r := GetRouter()
	n := negroni.Classic()
	n.UseHandler(r)
	http.ListenAndServe(":80", n)
}

// IndexHandler : Index route handler
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	resMap := make(map[string]interface{})
	resp, _ := MakeJSONResponse("I am the auctioner, sup dude!", resMap, true)
	SendJSONHttpResponse(w, resp, http.StatusOK)
}

func registerRoutes(router *mux.Router) {
	router.Methods("GET").Path("/").HandlerFunc(IndexHandler)
	router.Methods("POST").Path("/adplacement").HandlerFunc(AdPlacementHandler)
}

func GetRouter() *mux.Router {
	r := mux.NewRouter()
	registerRoutes(r)
	return r
}
