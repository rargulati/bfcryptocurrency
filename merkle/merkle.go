package main

import "hash"

type MerkleTree struct {
	root hash.Hash
	node *Node
}

type Node struct {
	Leaf  bool
	Block hash.Hash
}

func main() {

}

// plan of attack:
// 1. Read up on merkle trees, understand the structure of this thing
// 2.
