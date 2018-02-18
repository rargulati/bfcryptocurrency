# Merkle tree exercise

* Generate a Merkle tree given the following body of data, using SHA-2 as your hashing algorithm
  - Your data is the following blocks:
  - "We", "hold", "these", "truths", "to", "be, "self-evident", "that"
* In the internal stages, concatenate blocks as so "#{blockA.hash}||#{blockB.hash}"
* Do not use any special padding for leaf vs internal nodes
* The merkle root should be equal to `c4f66b2f97c9fb2fcb58b08b4f260d396b5c972ff4948c7deccc81fa34db1a44`

EDIT: with sha2, correct root hash is actually `4a359c93d6b6c9beaa3fe8d8e68935aa5b5081bd2603549af88dee298fbfdd0a`

Bonus: create a padding scheme so that arbitrary numbers of blocks can be Merkleized.

Bonus 2: add different padding to the leaves as opposed to internal nodes, so that preimage attacks are impossible.

Bonus 3: implement an interface for Merkle proofs. Have a `prove_inclusion(block)` function and a `verify_inclusion(proof, merkle_root)` function.
