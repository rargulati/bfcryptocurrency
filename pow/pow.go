package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type State struct {
	seen map[string]bool // seen tokens, token -> bool
	mu   sync.Mutex
}

// TODO: unsigned vs signed int for workFactor
func (s *State) mint(challenge int64, workFactor int) (int64, string) {
	token := int64(0)
outer:
	for {
		output := computeSingle(challenge, token)
		for _, byte := range output[:workFactor] {
			if byte != '0' {
				token++
				continue outer
			}
		}

		return token, output
	}
}

func computeSingle(challenge int64, token int64) string {
	ch := strconv.FormatInt(challenge, 10)
	concat := fmt.Sprintf("%s||%s", ch, strconv.FormatInt(token, 10))

	h := sha256.New()
	if _, err := h.Write([]byte(concat)); err != nil {
		panic(err)
	}
	output := fmt.Sprintf("%x", h.Sum(nil))
	return output
}

func (s *State) verify(challenge int64, workFactor int, token int64) bool {
	output := computeSingle(challenge, token)

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.seen[output] {
		return false
	}
	s.seen[output] = true

	zeros := strings.Repeat("0", workFactor)
	return strings.EqualFold(string(output[:workFactor]), zeros)
}

func main() {
	s := &State{
		seen: make(map[string]bool),
	}
	workFactor := 4
	for c := 0; c < 500; c++ {
		token, output := s.mint(int64(c), workFactor)
		if ok := s.verify(int64(c), workFactor, token); ok {
			fmt.Printf("Success at token: %d, output: %s\n", token, output)
		} else {
			fmt.Printf("Double Spend at token: %d, output: %s\n", token, output)
		}
	}
}
