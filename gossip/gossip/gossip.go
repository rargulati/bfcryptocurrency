package gossip

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Book string

type BookState struct {
	Favorite Book
	mu       sync.Mutex //guards Favorite and Version
	Books    []Book
	// initOnce sync.Once //ensure list of books grabbed only once
}

type Node struct {
	Port   string
	Router *mux.Router
	// transport http.RoundTripper
	State *BookState

	PeerList []string
	peerLock sync.Mutex

	RecentMessages map[string][]*Message // port -> Message
	messageLock    sync.Mutex

	ctx context.Context
}

type Message struct {
	// ID       uuid.UUID
	Port string
	// TTL      time.Duration
	Favorite Book
}

func NewNode(port int) *Node {
	p := strconv.Itoa(port)
	n := &Node{
		Port: p,
		State: &BookState{
			Books: make([]Book, 0),
		},
		// transport: &http.Transport{
		// 	IdleConnTimeout: 1 * time.Second,
		// 	DialContext: (&net.Dialer{
		// 		Timeout:   5 * time.Second,
		// 		KeepAlive: 5 * time.Second,
		// 		DualStack: true,
		// 	}).DialContext,
		// 	MaxIdleConns:          100,
		// 	ExpectContinueTimeout: 1 * time.Second,
		// },
		PeerList:       make([]string, 0),
		RecentMessages: make(map[string][]*Message),
		ctx:            context.Background(),
	}

	n.init()

	return n
}

func (n *Node) init() {
	var err error
	n.State.Books, err = getBooks("books_txt.txt")
	if err != nil {
		panic("Failed to parse books")
	}

	// racy
	go n.bootstrap()
	go n.resampleFavoriteBook()
	go n.broadcastStateLoop()
	go n.getPeersLoop()
	go n.dumpState()
}

func (n *Node) dumpState() {
	t := time.NewTimer(1)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			n.peerLock.Lock()
			log.Printf("Current PeerList: %#v\n", n.PeerList)
			n.peerLock.Unlock()

			n.messageLock.Lock()
			log.Printf("RecentMessages dump: %#v\n", n.RecentMessages)
			n.messageLock.Unlock()

			t.Reset(5 * time.Second)
		case <-n.ctx.Done():
			return
		}
	}
}

func (n *Node) broadcastStateLoop() {
	t := time.NewTimer(1)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			n.broadcastState()
			t.Reset(5 * time.Second)
		case <-n.ctx.Done():
			return
		}
	}
}

func (n *Node) broadcastState() {
	// get snapshot of current peers
	n.peerLock.Lock()
	peers := make([]string, len(n.PeerList))
	copy(peers, n.PeerList)
	n.peerLock.Unlock()

	// send gossip to all peers
	n.messageLock.Lock()
	defer n.messageLock.Unlock()

	for _, peer := range peers {
		// avoid sending ourself our own messages?
		if n.Port == peer {
			continue
		}
		for _, messages := range n.RecentMessages {
			// get most recent message for peer
			if len(messages) < 1 {
				continue
			}
			message := messages[len(messages)-1]

			msg, err := json.Marshal(struct {
				Message *Message
			}{
				Message: message,
			})
			if err != nil {
				panic(fmt.Sprintf("json error: %s", err))
			}

			peerLocation := fmt.Sprintf("localhost:%s", peer)
			url := fmt.Sprintf("http://%s/gossip/", peerLocation)
			// send our messages over the wire to this candidate
			// req, err := http.NewRequest("POST", url, bytes.NewBuffer(msg))
			// if err != nil {
			// 	panic(fmt.Sprintf("response creation error: %s", err))
			// }

			// res, err := n.transport.RoundTrip(req)
			res, err := http.Post(url, "application/json", bytes.NewBuffer(msg))
			if err != nil {
				log.Println(fmt.Sprintf("transport error: %s", err))
				continue
			}
			defer res.Body.Close()

			if code := res.StatusCode; code != 200 {
				b, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Println(fmt.Sprintf("couldn't read body: %v", err))
					continue
				}
				log.Println(fmt.Sprintf("not 200: %d and %#v", res.StatusCode, string(b)))
				continue
			}
			log.Println("sent successfull message %#v to member: %s", message, peer)
		}
	}
}

func (n *Node) GetPeerList() []string {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	return n.PeerList
}

func (n *Node) getPeersLoop() {
	t := time.NewTimer(1)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			n.pollAllPeers()
			t.Reset(5 * time.Second)
		case <-n.ctx.Done():
			return
		}
	}
}

func (n *Node) pollAllPeers() {
	// get snapshot of current peers
	n.peerLock.Lock()
	peers := make([]string, len(n.PeerList))
	copy(peers, n.PeerList)
	n.peerLock.Unlock()

	// get an idea of what peers others have
	otherPeers := make([]string, 0)
	for _, peer := range peers {
		// No need to get our own peerlist
		if peer == n.Port {
			continue
		}
		peerLocation := fmt.Sprintf("localhost:%s", peer)
		url := fmt.Sprintf("http://%s/peers", peerLocation)

		// req, err := http.NewRequest("GET", url, nil)
		// if err != nil {
		// 	panic(fmt.Sprintf("response creation error: %s", err))
		// }
		// res, err := n.transport.RoundTrip(req)
		res, err := http.Get(url)
		if err != nil {
			log.Println(fmt.Sprintf("transport error: %s", err))
			continue
		}
		defer res.Body.Close()

		if code := res.StatusCode; code != 200 {
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Println(fmt.Sprintf("couldn't read body: %v", err))
				continue
			}
			log.Println(fmt.Sprintf("not 200: %d and %#v", res.StatusCode, string(b)))
			continue
		}

		var externalPeers []string
		if err := json.NewDecoder(res.Body).Decode(&externalPeers); err != nil {
			panic(fmt.Sprintf("failed to unmarshal peers body with err: ", err))
		}

		for _, external := range externalPeers {
			otherPeers = append(otherPeers, external)
		}
	}

	seen := make(map[string]bool)
	// Add the ones we already know about
	for _, peer := range peers {
		seen[peer] = true
	}

	// Add the ones we don't know about, accounting for duplicates
	for _, peer := range otherPeers {
		if _, ok := seen[peer]; !ok {
			seen[peer] = true

			// update our peer list
			n.peerLock.Lock()
			n.PeerList = append(n.PeerList, peer)
			n.peerLock.Unlock()
		}
	}
}

func (n *Node) ProcessIncomingGossip(m *Message) error {
	n.messageLock.Lock()
	defer n.messageLock.Unlock()

	log.Printf("Process message %#v\n", m)
	if _, ok := n.RecentMessages[m.Port]; !ok {
		n.RecentMessages[m.Port] = make([]*Message, 0)
		n.RecentMessages[m.Port] = append(n.RecentMessages[m.Port], m)

		// add the peer to our known peers
		n.peerLock.Lock()
		n.PeerList = append(n.PeerList, m.Port)
		n.peerLock.Unlock()

		log.Printf("Added UNKNOWN peer %s and message %#v\n", m.Port, m)
		return nil
	}

	n.RecentMessages[m.Port] = append(n.RecentMessages[m.Port], m)
	log.Printf("Added KNOWN peer %s and message %#v\n", m.Port, m)
	return nil
}

func (n *Node) resampleFavoriteBook() {
	// first tick is immediate
	t := time.NewTimer(1)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			n.selectFavoriteBook()
			n.createNewBroadcastMessage()
			t.Reset(10 * time.Second)
		case <-n.ctx.Done():
			return
		}
	}
}

func (n *Node) selectFavoriteBook() {
	r := rand.Intn(len(n.State.Books)-1) + 1

	n.State.mu.Lock()
	n.State.Favorite = n.State.Books[r]
	n.State.mu.Unlock()
}

func (n *Node) createNewBroadcastMessage() {
	n.State.mu.Lock()
	favorite := n.State.Favorite
	n.State.mu.Unlock()

	n.messageLock.Lock()
	newMessage := &Message{Port: n.Port, Favorite: favorite}
	n.RecentMessages[n.Port] = append(n.RecentMessages[n.Port], newMessage)
	n.messageLock.Unlock()
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

func (n *Node) bootstrap() {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	// add yourself to your peerlist
	n.PeerList = append(n.PeerList, n.Port)

	if n.Port == "2702" {
		return
	}

	// add the bootstrap node to the peerlist
	n.PeerList = append(n.PeerList, "2702")

	// say hello to bootstrap node
	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)

		n.messageLock.Lock()
		messages, ok := n.RecentMessages[n.Port]
		if !ok {
			log.Println("Failed bootstrap peer exchange for, no message: ", n.Port)
			continue
		}
		message := messages[len(messages)-1]
		n.messageLock.Unlock()

		log.Printf("BOOTSTRAPPING and sending: %#v\n", message)

		msg, err := json.Marshal(struct {
			Message *Message
		}{
			Message: message,
		})
		if err != nil {
			panic(fmt.Sprintf("json error: %s", err))
		}

		url := fmt.Sprintf("http://%s/gossip/", "localhost:2702")
		// send our messages over the wire to this candidate
		// req, err := http.NewRequest("POST", url, bytes.NewBuffer(msg))
		// if err != nil {
		// 	panic(fmt.Sprintf("response creation error: %s", err))
		// }
		res, err := http.Post(url, "application/json", bytes.NewBuffer(msg))
		// res, err := n.transport.RoundTrip(req)
		if err != nil {
			log.Println(fmt.Sprintf("transport error: %s", err))
			continue
		}
		defer res.Body.Close()

		if code := res.StatusCode; code != 200 {
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Println(fmt.Sprintf("couldn't read body: %v", err))
				continue
			}
			log.Println(fmt.Sprintf("not 200: %d and %#v", res.StatusCode, string(b)))
			continue
		}
		return
	}
	return
}
