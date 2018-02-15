package node

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type Book string

type BookState struct {
	Favorite *Book
	mu       sync.Mutex //guards Favorite and Version
	Books    []Book
	initOnce sync.Once //ensure list of books grabbed only once
}

type BroadcastQueue struct {
	Messages []*Message
	mu       sync.Mutex
}

type MemberList struct {
	nodes   []*Node
	nodeMap map[string]*Node // Addr.String() -> Node
	mu      sync.Mutex
}

type Node struct {
	Addr      net.IP
	Port      string
	Router    *mux.Router
	transport http.RoundTripper

	// Queue *BroadcastQueue
	cache *MemberList
	State *BookState

	tracker     updateTracker
	trackerLock sync.Mutex

	incomingMessage chan *Message
	outgoingMessage chan *Message

	ctx context.Context
}

type Message struct {
	ID      uuid.UUID
	Port    string
	Version update
	// TTL      time.Duration
	Favorite *Book
}

func NewNode() *Node {
	n := &Node{
		State: &BookState{
			Books: make([]Book, 0),
		},
		cache: &MemberList{
			nodes:   make([]*Node, 0),
			nodeMap: make(map[string]*Node, 0),
		},
		// Queue: &BroadcastQueue{
		// 	Messages: make([]*Message, 0),
		// },
		tracker: updateTracker{
			current: 0,
			seen:    make(map[update]bool),
		},
		incomingMessage: make(chan *Message, 1),
		outgoingMessage: make(chan *Message, 1),
		ctx:             context.Background(),
	}
	n.init()

	return n
}

func (n *Node) jointHostPort(host string, port string) string {
	return net.JoinHostPort(host, port)
}

func (n *Node) init() {
	n.State.initOnce.Do(func() {
		var err error
		n.State.Books, err = getBooks("books_txt.txt")
		if err != nil {
			panic("Failed to parse books")
		}
		go n.resampleFavoriteBook()
		go n.acceptIncoming()
		go n.sendStateLoop()
	})
}

func (n *Node) acceptIncomingLoop() {
	t := time.NewTimer(1)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			n.acceptIncoming()
			t.Reset(100 * time.Millisecond)
		case <-n.ctx.Done():
			return
		}
	}
}

func (n *Node) acceptIncoming() {
	for message := range n.incomingMessage {
		isNew := n.tracker.See(message.Version)
		if isNew {
			n.outgoingMessage <- message
		}
	}
}

func (n *Node) ProcessIncomingGossip(m *Message) error {
	// debug
	// fmt.Println("Message: %+v", m)
	// fmt.Println("Message: %#v", m)
	// TODO: do error type checking
	n.incomingMessage <- m
	return nil
}

func (n *Node) sendStateLoop() {
	t := time.NewTimer(1)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			n.sendState()
			t.Reset(100 * time.Millisecond)
		case <-n.ctx.Done():
			return
		}
	}
}

func (n *Node) sendState() {
	// TODO: choose a random node?
	n.cache.mu.Lock()
	nodes := n.cache.nodes
	n.cache.mu.Unlock()

	for outgoing := range n.outgoingMessage {
		for _, node := range nodes {
			url := fmt.Sprintf("%s/gossip/", node.jointHostPort(node.Addr.String(), node.Port))
			msg, err := json.Marshal(struct {
				Message *Message
			}{
				Message: outgoing,
			})
			if err != nil {
				panic("resp creation error")
			}

			// send our messages over the wire to this candidate
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(msg))
			if err != nil {
				panic("resp creation error")
			}
			res, err := n.transport.RoundTrip(req)
			if err != nil {
				panic("transport errror")
			}
			defer res.Body.Close()

			if code := res.StatusCode; code != 200 {
				panic("Not 200")
			}
		}
	}
}

func randNode(nodes []*Node) *Node {
	if len(nodes) == 0 {
		return nil
	}
	// idx := rand.Intn(len(nodes)-1) + 1 // get a random number
	// idx := randomOffset(len(nodes)) // find our random offset
	idx := rand.Int31n(int32(len(nodes) - 1))
	// fmt.Println(idx)
	return nodes[idx] // return our random node
}

func randomOffset(n int) int {
	return int(rand.Uint32() % uint32(n))
}

// TODO: pass in current sequence number
// func (n *Node) pendingBroadcasts() []*Message {
// 	n.Queue.mu.Lock()
// 	bcasts := n.Queue.Messages
// 	msgs := make([]*Message, len(bcasts))
// 	copy(msgs, bcasts)
// 	n.Queue.mu.Unlock()

// 	return msgs
// }

func (n *Node) createNewBroadcastMessage(b *Book, p string) *Message {
	id, err := uuid.NewV4()
	if err != nil {
		panic("error generating new uuid")
	}

	n.trackerLock.Lock()
	current := n.tracker.current
	n.tracker.current++
	n.trackerLock.Unlock()

	m := &Message{
		ID:      id,
		Port:    p,
		Version: current,
		// TTL:      time.Duration(1 * time.Second),
		Favorite: b,
	}
	return m
}

func (n *Node) resampleFavoriteBook() {
	t := time.NewTimer(1)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			n.selectFavoriteBook()
			t.Reset(10 * time.Second)
		case <-n.ctx.Done():
			return
		}
	}
}

func (n *Node) selectFavoriteBook() {
	r := rand.Intn(len(n.State.Books)-1) + 1
	n.State.mu.Lock()
	n.State.Favorite = &n.State.Books[r]
	// n.State.Version = n.State.Version + 1

	// create a new message
	m := n.createNewBroadcastMessage(n.State.Favorite, n.Port)
	n.State.mu.Unlock()

	// add it to the broadcast queue
	n.incomingMessage <- m
	// n.Queue.mu.Lock()
	// n.Queue.Messages = append(n.Queue.Messages, m)
	// n.Queue.mu.Unlock()
}

func getBooks(path string) ([]Book, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	books := make([]Book, 0)
	for scanner.Scan() {
		str := scanner.Text()
		books = append(books, Book(str))
	}
	return books, nil
}
