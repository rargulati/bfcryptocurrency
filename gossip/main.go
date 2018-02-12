package main

import (
	"bufio"
	"math/rand"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

// Node state
type Node struct {
	Router *mux.Router
	// Peer
	// favorite book
	// List of books that we know about
	// Protocol
	bs *BookState
}

type BookState struct {
	Favorite *Book
	mu       sync.Mutex //guards Book
	Books    []Book
	initOnce sync.Once //ensure list of books grabbed only once
}

type Book string

func main() {
	// initialize application
	// select favorite book at random in a goroutine
	n := NewNode()
	n.bs.initOnce.Do(func() {
		var err error
		n.bs.Books, err = getBooks("book_list.txt")
		if err != nil {
			panic("Failed to parse books")
		}
		go n.resampleFavoriteBook()
	})
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) resampleFavoriteBook() {
	r := rand.Intn(len(n.bs.Books) - 0)
	n.bs.mu.Lock()
	n.bs.Favorite = &n.bs.Books[r]
	n.bs.mu.Unlock()
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
