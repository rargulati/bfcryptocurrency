package main

import (
	"crypto/sha256"
	"fmt"
	"hash"
)

type Node struct {
	left   *Node
	right  *Node
	parent *Node
	hash   string
	data   string
}

type MerkleTree struct {
	root       *Node
	merkleRoot hash.Hash
	leafs      []*Node
}

func hashData(data string) string {
	// h := sha256.New()
	// _, _ = h.Write([]byte(data))
	sum := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", sum)
}

func BuildSubTrees(leafs []*Node) (*Node, error) {
	var nodes []*Node
	// for each leaf node (that's every two)
	for i := 0; i < len(leafs); i += 2 {
		// intermediate := append([]byte(leafs[i].hash), []byte(leafs[i+1].hash)...)
		intermediate := fmt.Sprintf("%s%s", leafs[i].hash, leafs[i+1].hash)
		h := hashData(intermediate)
		fmt.Printf("HASH of intermediate(%d): %s\n", i, h)
		n := &Node{
			left:  leafs[i],
			right: leafs[i+1],
			hash:  h,
		}
		nodes = append(nodes, n)
		// set our parent node for our new sub-tree nodes
		leafs[i].parent = n
		leafs[i+1].parent = n
		// We're working at the last layer, we're the root
		if len(leafs) == 2 {
			return n, nil
		}
	}
	fmt.Println(len(nodes))
	// recursively call this with our new layer
	return BuildSubTrees(nodes)
}

func NewMerkleTree(data []string) *MerkleTree {
	m := &MerkleTree{}
	// create and populate the leaf layer of the merkle tree
	for _, s := range data {
		// TODO: avoid pre-image attacks with padding our data
		n := &Node{
			data: s,
			hash: hashData(s),
		}
		m.leafs = append(m.leafs, n)
	}
	// odd case - where we have an odd number of leaves
	// pad with an empty string
	if len(m.leafs)%2 == 1 {
		m.leafs = append(m.leafs, &Node{data: "", hash: hashData("")})
	}
	return m
}

func main() {
	data := []string{"We", "hold", "these", "truths", "to", "be", "self-evident", "that"}
	m := NewMerkleTree(data)
	root, err := BuildSubTrees(m.leafs)
	if err != nil {
		fmt.Println("err: ", err)
	}
	fmt.Printf("ROOT: %s", root.hash)
}
