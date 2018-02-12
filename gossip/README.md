# Building a gossip protocol

## Goal
We're going to build a gossip protocol that gossips with a network of peers about their current favorite books. Each node will connect to only one other node, and all nodes should know the favorite books of all other nodes in the system. We'll have nodes run as separate processes on different ports on your local machine, and they'll pass JSON-encoded messages via HTTP (to make our lives easier).

## State
We'll use a gossip protocol to keep track of each node's current favorite book. You can find a list of books in `books.txt`.

Each node's book should be randomly re-sampled from the pool of all books once every ~10 seconds. Once it chooses a new favorite book, it should flood its peers with this message.

You need to have each node keep track of their own incrementing version number, so we can keep track of their state and order messages. In a gossip protocol we will often receive messages out of order, so we need to know which one is most recent.

The node should also keep a cache of the recent messages it's received. Normally we'd want to cull this, but for now we can just let it grow in memory.

## Endpoints
For simplicity, we'll bootstrap the network with one other node. That argument will be passed via the command line.

Each node needs the endpoints:

* GET /peers/ (for bootstrapping into the network)
* POST /gossip/ (for sending gossip between nodes)

## Message format
Your messages will need the following:

* UUID (for deduplication)
* Originating port (your identity)
* Version number
* TTL
* Payload
