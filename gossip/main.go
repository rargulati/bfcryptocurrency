package main

import (
	"flag"
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
	// Parse options from the command line
	listenF := flag.Int("l", 0, "wait for incoming connections")
	flag.Parse()
	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	// TODO: command line setttings ie file path
	n := node.NewNode()
	r := mux.NewRouter()
	n.Router = r
	n.Router.HandleFunc("/", IndexHandler)
	n.Router.HandleFunc("/peers", n.PeersHandler).Methods("GET")
	n.Router.HandleFunc("/gossip", n.GossipHandler).Methods("POST")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe("localhost:8000", n.Router))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}
