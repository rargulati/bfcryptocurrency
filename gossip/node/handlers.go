package node

import (
	"encoding/json"
	"net/http"
)

func (n *Node) PeersHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Peers!\n"))
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

	// fin
	w.Write([]byte("Gossip!\n"))
}
