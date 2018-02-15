package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rargulati/blockchain-class/gossip/node"
)

func init() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// TODO: command line setttings ie file path
	n := node.NewNode()
	r := mux.NewRouter()
	n.Router = r
	n.Router.HandleFunc("/", YourHandler)
	n.Router.HandleFunc("/peers", n.PeersHandler).Methods("GET")
	n.Router.HandleFunc("/gossip", n.GossipHandler).Methods("POST")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe("localhost:8000", n.Router))
}

func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}
