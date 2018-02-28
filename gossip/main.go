package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	node "github.com/rargulati/blockchain-class/gossip/gossip"
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
	n := node.NewNode(*listenF)
	r := mux.NewRouter()
	n.Router = r
	n.Router.HandleFunc("/", IndexHandler)
	n.Router.HandleFunc("/peers", n.PeersHandler).Methods("GET")
	n.Router.HandleFunc("/gossip/", n.GossipHandler).Methods("POST")

	// Bind to a port and pass our router in
	lis := fmt.Sprintf(":%d", *listenF)
	log.Printf("Listening at: %s\n", lis)
	log.Fatal(http.ListenAndServe(lis, n.Router))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}
