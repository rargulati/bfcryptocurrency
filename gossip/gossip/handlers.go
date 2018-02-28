package gossip

import (
	"encoding/json"
	"log"
	"net/http"
)

func (n *Node) PeersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Allow", "OPTIONS,GET")
	if r.Method == "OPTIONS" {
		log.Println("Options")
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Invalid Method", 400)
		return
	}

	// Get our peer list
	peers := n.GetPeerList()
	// log.Printf("PeerList: %#v", peers)

	if err := json.NewEncoder(w).Encode(&peers); err != nil {
		http.Error(w, "Invalid Server Error", 500)
		return
	}
	return
}

func (n *Node) GossipHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Allow", "OPTIONS,POST")
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Invalid Method", 400)
		return
	}

	var v struct {
		Message *Message
	}
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	// pass off the gossip to our node
	if err := n.ProcessIncomingGossip(v.Message); err != nil {
		http.Error(w, "Invalid Message", 400)
		return
	}
	// log.Printf("GOSSIP: %#v", v.Message)
	return
}
