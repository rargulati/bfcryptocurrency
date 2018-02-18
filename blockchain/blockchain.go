package main

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

const workFactor = 4

type Block struct {
	content      string
	previousHash string
	nonce        string
}

type Blockchain struct {
	blocks []*Block
}

func NewBlock(data string) *Block {
	return &Block{}
}

func NewBlockchain(data []string) *Blockchain {
	bs := make([]*Block, len(data))
	for i, d := range data {
		b := &Block{content: d}
		if i > 0 {
			// b.previousHash =
		}
		bs = append(bs, b)
	}
	return &Blockchain{blocks: bs}
}

func proofOfWork(data string) (string, bool) {
	h := sha256.New()
	if _, err := h.Write([]byte(data)); err != nil {
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

func main() {
	b := NewBlockchain(strings.Split("this is a string", " "))
	fmt.Println("Printing Blockchain...")
	for _, block := range b.blocks {
		fmt.Printf("Content: %s, PreviousHash: %s, Nonce: %d\n", block.content, block.previousHash, block.nonce)
	}
}
