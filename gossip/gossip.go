package main

import (
	"sync"
	"time"

	"github.com/satori/go.uuid"
)

type Message struct {
	ID      uuid.UUID
	Port    string
	Version string
	TTL     time.Duration
	Payload string
}

type memberState struct {
	Port        string
	Protocol    string
	MemberState int
}

type Peer struct {
	// A cache of peers (port) and their current favorite book
	cache map[string]string
	mu    sync.Mutex
}

func init() {
	// Seed the random number generator
	// rand.Seed(time.Now().UnixNano())
	// load the books into memory
	// getlist("books.txt")
}
