package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"strings"

	"github.com/agl/ed25519"
)

const workFactor = 4

type Block struct {
	content      string
	previousHash string
	hash         string
	nonce        int64
	tx           *Transaction
}

type Blockchain struct {
	blocks []*Block
}

type Transaction struct {
	from      *Ed25519PublicKey
	to        *Ed25519PublicKey
	amount    int
	signature string
}

type Ed25519PrivateKey struct {
	sk *[64]byte
	pk *[32]byte
}

type Ed25519PublicKey struct {
	k *[32]byte
}

func (k *Ed25519PrivateKey) Sign(msg []byte) ([]byte, error) {
	out := ed25519.Sign(k.sk, msg)
	return (*out)[:], nil
}

func (k *Ed25519PublicKey) Verify(data []byte, sig []byte) (bool, error) {
	var asig [64]byte
	copy(asig[:], sig)
	return ed25519.Verify(k.k, data, &asig), nil
}

func NewBlockchain() *Blockchain {
	return &Blockchain{blocks: make([]*Block, 0)}
}

func (bc *Blockchain) InitializeGenesisBlock(d string) {
	b := &Block{content: d}
	b.previousHash = "Genesis"
	blockHash := hashBlock(fmt.Sprintf("%s||%s", d, b.previousHash))
	b.nonce, b.hash = proofOfWork(blockHash, workFactor)
	bc.blocks = append(bc.blocks, b)
	return
}

func (bc *Blockchain) AddBlocks(data []string) error {
	if bc.blocks == nil {
		return fmt.Errorf("Must first initialize a genesis block")
	}

	for i, d := range data {
		b := &Block{content: d}
		var pb *Block
		if i == 0 {
			pb = bc.blocks[0]
		} else {
			pb = bc.blocks[i-1]
		}

		b.previousHash = pb.hash
		concat := fmt.Sprintf("%s||%s", d, b.previousHash)
		b.nonce, b.hash = proofOfWork(hashBlock(concat), workFactor)

		bc.blocks = append(bc.blocks, b)
	}
	return nil
}

func hashBlock(data string) string {
	h := sha256.New()
	if _, err := h.Write([]byte(data)); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func proofOfWork(challenge string, workFactor int) (int64, string) {
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

func computeSingle(challenge string, token int64) string {
	concat := fmt.Sprintf("%s||%s", challenge, strconv.FormatInt(token, 10))

	h := sha256.New()
	if _, err := h.Write([]byte(concat)); err != nil {
		panic(err)
	}
	output := fmt.Sprintf("%x", h.Sum(nil))
	return output
}

func verify(challenge string, workFactor int, token int64) bool {
	output := computeSingle(challenge, token)
	zeros := strings.Repeat("0", workFactor)
	return strings.EqualFold(string(output[:workFactor]), zeros)
}

func main() {
	data := strings.Split("this is a string", " ")
	b := NewBlockchain()
	b.InitializeGenesisBlock(data[0])
	b.AddBlocks(data[1:])
	fmt.Println("Printing Blockchain...")
	for _, block := range b.blocks {
		fmt.Printf("Content: %s, Hash: %s, PreviousHash: %s, Nonce: %d\n", block.content, block.hash, block.previousHash, block.nonce)
	}
}
