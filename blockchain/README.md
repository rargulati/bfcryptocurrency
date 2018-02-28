# Build a blockchain!

In this assignment, you'll be building a blockchain from scratch. You'll want to build a block class, and then a blockchain class that includes all blocks.

* Block
  - `content` (string)
  - `previous_hash` (string)
  - `current_hash` (string)
  - `nonce` (string/number)

* Blockchain
  - `blocks` (array)

I should be able to initialize your blockchain as `Blockchain.new("this is a string".split(" "))` and have it automatically create blocks, compute proof of work for each block, and chain them all together in a blockchain. You should also be able to append new blocks to this blockchain.

Note: if you finish early and have some extra time, I’d encourage you to create a Transaction class for your blockchain. Transactions should have a `from` (public key), `to` (public key), `amount` (int), and `signature` (string). You should make your blockchain validate all transaction signatures, and ensure that no balances go into the negative.
