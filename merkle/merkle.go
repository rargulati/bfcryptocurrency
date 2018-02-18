package main

import (
	"crypto/sha256"
	"fmt"
	"hash"
)

type Node struct {
	Left   *Node
	Right  *Node
	Parent *Node
	hash   string
	data   string
}

type MerkleTree struct {
	Root       *Node
	merkleRoot hash.Hash
	leafs      []*Node
}

func hashData(data ...[]byte) string {
	h := sha256.New()
	for _, d := range data {
		_, _ = h.Write(d)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func BuildSubTrees(leafs []*Node) (*Node, error) {
	// if m == nil {
	// 	return nil, fmt.Errorf("Cannot insert value into nil merkle tree")
	// }
	// for each leaf node (that's every two)
	var nodes []*Node
	for i := 0; i < len(leafs); i += 2 {
		intermediate := append([]byte(leafs[i].hash), []byte(leafs[i+1].hash)...)
		h := hashData([]byte(intermediate))
		fmt.Printf("HASH OF intermediate(%d): %s\n", i, h)
		n := &Node{
			Left:  leafs[i],
			Right: leafs[i+1],
			hash:  h,
		}
		nodes = append(nodes, n)
		// set our parent node for our new sub-tree nodes
		leafs[i].Parent = n
		leafs[i+1].Parent = n
		// If we only have two nodes left, we're the root
		if len(leafs) == 2 {
			return n, nil
		}
	}
	fmt.Println(len(nodes))
	return BuildSubTrees(nodes)
}

func NewMerkleTree(data []string) *MerkleTree {
	m := &MerkleTree{}
	// create and populate the leaf layer of the merkle tree
	for _, s := range data {
		n := &Node{
			data: s,
			hash: hashData([]byte(s)),
		}
		m.leafs = append(m.leafs, n)
	}
	// odd case - where we have an odd number of leaves
	// if len(m.leafs)%2 == 1 {
	// 	m.leafs = append(m.leafs, &Node{data: "", hash: hashData([]byte(""))})
	// }
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
