package gossip

import "testing"

func TestGossipsubTreeTopology(t *testing.T) {
	// 	ctx, cancel := context.WithCancel(context.Background())
	// 	defer cancel()

	// 	hosts := getNetHosts(t, ctx, 10)
	// 	psubs := getGossipsubs(ctx, hosts)

	// 	connect(t, hosts[0], hosts[1])
	// 	connect(t, hosts[1], hosts[2])
	// 	connect(t, hosts[1], hosts[4])
	// 	connect(t, hosts[2], hosts[3])
	// 	connect(t, hosts[0], hosts[5])
	// 	connect(t, hosts[5], hosts[6])
	// 	connect(t, hosts[5], hosts[8])
	// 	connect(t, hosts[6], hosts[7])
	// 	connect(t, hosts[8], hosts[9])

	// 	/*
	// 		[0] -> [1] -> [2] -> [3]
	// 		 |      L->[4]
	// 		 v
	// 		[5] -> [6] -> [7]
	// 		 |
	// 		 v
	// 		[8] -> [9]
	// 	*/

	// 	var chs []*Subscription
	// 	for _, ps := range psubs {
	// 		ch, err := ps.Subscribe("fizzbuzz")
	// 		if err != nil {
	// 			t.Fatal(err)
	// 		}

	// 		chs = append(chs, ch)
	// 	}

	// 	// wait for heartbeats to build mesh
	// 	time.Sleep(time.Second * 2)

	// 	assertPeerLists(t, hosts, psubs[0], 1, 5)
	// 	assertPeerLists(t, hosts, psubs[1], 0, 2, 4)
	// 	assertPeerLists(t, hosts, psubs[2], 1, 3)

	// 	checkMessageRouting(t, "fizzbuzz", []*PubSub{psubs[9], psubs[3]}, chs)
}

// func getNetHosts(t *testing.T, ctx context.Context, n int) []host.Host {
// 	var out []host.Host

// 	for i := 0; i < n; i++ {
// 		netw := netutil.GenSwarmNetwork(t, ctx)
// 		h := bhost.NewBlankHost(netw)
// 		out = append(out, h)
// 	}

// 	return out
// }

// func connect(t *testing.T, a, b host.Host) {
// 	pinfo := a.Peerstore().PeerInfo(a.ID())
// 	err := b.Connect(context.Background(), pinfo)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
// func assertPeerLists(t *testing.T, hosts []host.Host, ps *PubSub, has ...int) {
// 	peers := ps.ListPeers("")
// 	set := make(map[peer.ID]struct{})
// 	for _, p := range peers {
// 		set[p] = struct{}{}
// 	}

// 	for _, h := range has {
// 		if _, ok := set[hosts[h].ID()]; !ok {
// 			t.Fatal("expected to have connection to peer: ", h)
// 		}
// 	}
// }

// func checkMessageRouting(t *testing.T, topic string, pubs []*PubSub, subs []*Subscription) {
// 	data := make([]byte, 16)
// 	rand.Read(data)

// 	for _, p := range pubs {
// 		err := p.Publish(topic, data)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		for _, s := range subs {
// 			assertReceive(t, s, data)
// 		}
// 	}
// }
