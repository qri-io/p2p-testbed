package main

import (
	"context"
	"sync"

	logging "github.com/ipfs/go-log"
	// logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"

	opentracing "github.com/opentracing/opentracing-go"

	pstore "gx/ipfs/QmPgDWmTmuzvP7QE5zwo1TmjbJme9pmZHNujB2453jkCTr/go-libp2p-peerstore"
	testutil "gx/ipfs/QmWRCn8vruNAzHx8i6SAXinuheRitKEGu8c7m26stKvsYx/go-testutil"
	ma "gx/ipfs/QmXY77cVe7rVRQXZZQRioukUM7aRW3BTcAgJe12MCtb3Ji/go-multiaddr"
	peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
)

const testbedProtocolID = "/testbed"

// Setup creates & connects a new network
func Setup(ctx context.Context, peerCount int) (peers []*Peer, err error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "Setup")
	span.SetTag("peerCount", peerCount)

	peers, err = NewNetwork(ctx, peerCount)
	if err != nil {
		return
	}

	err = connectPeers(ctx, peers)
	span.Finish()
	return
}

// NewNetwork creates a new network with num peers
func NewNetwork(ctx context.Context, num int) ([]*Peer, error) {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "NewNetwork")
	peers := make([]*Peer, 0, num)

	for i := 0; i < num; i++ {
		params, err := testutil.RandPeerNetParams()
		if err != nil {
			return nil, err
		}

		p, err := NewPeer(ctx, params)
		if err != nil {
			log.Errorf("error creating peer: %s", err.Error())
			return nil, err
		}

		peers = append(peers, p)
	}
	span.Finish()
	return peers, nil
}

// connectPeers connects a set of peers together
func connectPeers(ctx context.Context, peers []*Peer) error {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "connectPeers")

	var wg sync.WaitGroup
	connect := func(n *Peer, dst peer.ID, addr ma.Multiaddr) {
		log.Debugf("dialing %s from %s\n", n.ID, dst)

		evt := log.EventBegin(ctx, "DialPeer", logging.LoggableMap{
			"from": n.ID.String(),
			"to":   dst.String(),
		})

		n.Peerstore.AddAddr(dst, addr, pstore.PermanentAddrTTL)
		if _, err := n.Host.Network().DialPeer(ctx, dst); err != nil {
			log.Errorf("error swarm dialing to peer", err)
			return
		}

		evt.Done()
		wg.Done()
	}

	log.Infof("Connecting swarms simultaneously.")
	for i, s1 := range peers {
		for _, s2 := range peers[i+1:] {
			wg.Add(1)
			connect(s1, s2.Host.Network().LocalPeer(), s2.Host.Network().ListenAddresses()[0]) // try the first.
		}
	}
	wg.Wait()
	span.Finish()

	// for _, n := range peers {
	// 	log.Debugf("%s swarm routing table: %s\n", n.ID, n.Peerstore.Peers())
	// }
	return nil
}
