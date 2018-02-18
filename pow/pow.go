package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
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
		output, ok := computeSingleSha(challenge, workFactor, token)
		if !ok {
			token++
			continue outer
		}

		return token, output
	}
}

func computeSingleSha(challenge int64, workFactor int, token int64) (string, bool) {
	ch := strconv.FormatInt(challenge, 10)
	concat := fmt.Sprintf("%s%s", ch, strconv.FormatInt(token, 10))

	h := sha256.New()
	if _, err := h.Write([]byte(concat)); err != nil {
		panic(err)
	}
	output := fmt.Sprintf("%x", h.Sum(nil))
	for _, byte := range output[:workFactor] {
		if byte != '0' {
			return "", false
		}
	}
	return output, true
}

func (s *State) verify(challenge int64, workFactor int, token int64) bool {
	output, ok := computeSingleSha(challenge, workFactor, token)
	if !ok {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.seen[output] {
		return false
	}
	s.seen[output] = true
	return true
}

func main() {
	s := &State{
		seen: make(map[string]bool),
	}
	workFactor := 4
	for c := 0; c < 5000; c++ {
		token, output := s.mint(int64(c), workFactor)
		if ok := s.verify(int64(c), workFactor, token); ok {
			fmt.Printf("Success at token: %d, output: %s\n", token, output)
		} else {
			fmt.Printf("Failure at token: %d, output: %s\n", token, output)
		}
	}
}
