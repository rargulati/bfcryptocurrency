package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"strconv"
	"strings"

	crand "crypto/rand"

	"github.com/agl/ed25519"
)

const workFactor = 4

type Block struct {
	previousHash string
	hash         string
	nonce        int64
	tx           *Transaction
}

type Blockchain struct {
	blocks   []*Block
	accounts map[string]int // account balance
}

type Transaction struct {
	from      *Account
	to        *Account
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

type Account struct {
	name    string
	pubKey  *Ed25519PublicKey
	privKey *Ed25519PrivateKey
}

func NewKey(src io.Reader) (*Ed25519PrivateKey, *Ed25519PublicKey, error) {
	pub, priv, err := ed25519.GenerateKey(src)
	if err != nil {
		return nil, nil, err
	}
	return &Ed25519PrivateKey{sk: priv, pk: pub}, &Ed25519PublicKey{k: pub}, nil
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
	return &Blockchain{blocks: make([]*Block, 0), accounts: make(map[string]int)}
}

func (bc *Blockchain) updateAccountBalance(from, to *Account, value int) error {
	fromBalance, ok := bc.accounts[from.name]
	if !ok {
		return fmt.Errorf("No account with name %s exists", from.name)
	}
	if fromBalance-value < 0 {
		return fmt.Errorf("Can't send %d coins; will put %s in negative\n", value, from.name)
	}

	bc.accounts[from.name] = fromBalance - value
	if toBalance, ok := bc.accounts[to.name]; !ok {
		bc.accounts[to.name] = value
	} else {
		bc.accounts[to.name] = toBalance + value
	}
	return nil
}

func (bc *Blockchain) InitializeGenesisBlock(tx *Transaction) {
	bc.accounts[tx.from.name] = tx.amount
	b := &Block{tx: tx}
	b.previousHash = "Genesis"
	blockHash := hashBlock(fmt.Sprintf("%s||%s", tx.signature, b.previousHash))
	b.nonce, b.hash = proofOfWork(blockHash, workFactor)
	bc.updateAccountBalance(tx.from, tx.to, tx.amount)
	bc.blocks = append(bc.blocks, b)
	return
}

func (bc *Blockchain) AddBlocks(txs []*Transaction) error {
	if bc.blocks == nil {
		return fmt.Errorf("Must first initialize a genesis block")
	}

	for i, tx := range txs {
		b := &Block{tx: tx}
		var pb *Block
		if i == 0 {
			pb = bc.blocks[0]
		} else {
			pb = bc.blocks[i-1]
		}

		b.previousHash = pb.hash
		concat := fmt.Sprintf("%s||%s", tx.signature, b.previousHash)
		b.nonce, b.hash = proofOfWork(hashBlock(concat), workFactor)
		bc.updateAccountBalance(tx.from, tx.to, tx.amount)

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

func ValidateBlocks(bls []*Block) error {
	for i, bl := range bls {
		// verify pow
		concat := fmt.Sprintf("%s||%s", bl.tx.signature, bl.previousHash)
		if ok := verify(hashBlock(concat), workFactor, bl.nonce); !ok {
			return fmt.Errorf("Failed on block (%d): %#v\n", i, bl)
		}
		// verify sig
		ok, err := bl.tx.from.pubKey.Verify([]byte(strconv.Itoa(bl.tx.amount)), []byte(bl.tx.signature))
		if err != nil {
			return fmt.Errorf("Failed to verify with err: %s\n", err)
		}
		if !ok {
			return fmt.Errorf("Failed on block (%d): %#v\n", i, bl)
		}
	}
	return nil
}

func main() {
	accounts := []*Account{
		{name: "Jenn"},
		{name: "Radhika"},
		{name: "Raghav"},
		{name: "Rajni"},
		{name: "Ravi"},
	}
	for _, account := range accounts {
		priv, pub, err := NewKey(crand.Reader)
		if err != nil {
			panic(err)
		}
		account.pubKey = pub
		account.privKey = priv
	}

	txs := []*Transaction{
		{from: accounts[0], to: accounts[1], amount: 5},
		{from: accounts[1], to: accounts[2], amount: 4},
		{from: accounts[2], to: accounts[3], amount: 3},
		{from: accounts[3], to: accounts[4], amount: 2},
		{from: accounts[4], to: accounts[0], amount: 1},
	}

	for _, tx := range txs {
		sig, err := tx.from.privKey.Sign([]byte(strconv.Itoa(tx.amount)))
		if err != nil {
			panic(fmt.Sprintf("Signing err: %s\n", err))
		}
		tx.signature = string(sig)
	}

	b := NewBlockchain()
	b.InitializeGenesisBlock(txs[0])
	b.AddBlocks(txs[1:])
	fmt.Println("Validating Blockchain...")
	if err := ValidateBlocks(b.blocks); err != nil {
		fmt.Printf("Failed to ValidateBlocks with err: %s\n", err)
	}
	fmt.Println("Print out Blockchain accounts")

	for name, value := range b.accounts {
		fmt.Printf("Account %s with value %d\n", name, value)
	}

	fmt.Println("Blockchain Valid")
}
