# Build a blockchain!

In this assignment, you'll be building a blockchain from scratch. You'll want to build a block class, and then a blockchain class that includes all blocks.

* Block
  - `content` (string)
  - `previous_hash` (string)
  - `nonce` (string/number)

* Blockchain
  - `blocks` (array)

I should be able to initialize your blockchain as `Blockchain.new("this is a string".split(" "))` and have it automatically create blocks, compute proof of work for each block, and chain them all together in a blockchain. You should also be able to append new blocks to this blockchain.
