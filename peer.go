package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	yamux "gx/ipfs/QmNWCEvi7bPRcvqAV8AKLGVNoQdArWi7NJayka2SM4XtRe/go-smux-yamux"
	net "gx/ipfs/QmNa31VPzC561NWwRsJLE7nGYZYuuD2QfpK2b1q9BK54J1/go-libp2p-net"
	pstore "gx/ipfs/QmPgDWmTmuzvP7QE5zwo1TmjbJme9pmZHNujB2453jkCTr/go-libp2p-peerstore"
	msmux "gx/ipfs/QmVniQJkdzLZaZwzwMdd3dJTvWiJ1DQEkreVy6hs6h7Vk5/go-smux-multistream"
	testutil "gx/ipfs/QmWRCn8vruNAzHx8i6SAXinuheRitKEGu8c7m26stKvsYx/go-testutil"
	ma "gx/ipfs/QmXY77cVe7rVRQXZZQRioukUM7aRW3BTcAgJe12MCtb3Ji/go-multiaddr"
	peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
	host "gx/ipfs/Qmc1XhrFEiSeBNn3mpfg6gEuYCt5im2gYmNVmncsvmpeAk/go-libp2p-host"
	swarm "gx/ipfs/QmdQFrFnPrKRQtpeHKjZ3cVNwxmGKKS2TvhJTuN9C9yduh/go-libp2p-swarm"
	bhost "gx/ipfs/QmefgzMbKZYsmHFkLqxgaTBG9ypeEjrdWRD5WXH4j1cWDL/go-libp2p/p2p/host/basic"
)

// NewPeer creates a new peer from testutil params
func NewPeer(ctx context.Context, params *testutil.PeerNetParams) (peer *Peer, err error) {
	ps := pstore.NewPeerstore()
	if err = ps.AddPrivKey(params.ID, params.PrivKey); err != nil {
		return
	}
	if err = ps.AddPubKey(params.ID, params.PubKey); err != nil {
		return
	}

	host, err := makeBasicHost(ctx, params.ID, []ma.Multiaddr{params.Addr}, ps)
	if err != nil {
		return
	}

	state := &sync.Map{}
	state.Store("pet", RandomAnimal())

	peer = &Peer{
		ID:        params.ID,
		Host:      host,
		Peerstore: ps,
		ctx:       ctx,
		state:     state,
	}

	host.SetStreamHandler(testbedProtocolID, peer.TestbedHandler)
	return
}

// Peer is a peer in a (simulated) peer-2-peer network
type Peer struct {
	// ctx is the base operating context
	ctx context.Context
	// ID is the peer's identifier on the network
	ID peer.ID
	// Host carries all p2p protocols & services
	// for interacting on the network
	Host host.Host
	// Peerstore keeps a list
	Peerstore pstore.Peerstore
	// state is an in-memory map of any state
	state *sync.Map
}

// RandomAnimal generates a random animal, used for demonstration purposes
func RandomAnimal() string {
	animals := []string{
		"doggo",
		"cat",
		"giraffe",
		"lizard",
		"smurf",
		"chinchilla",
		"cryptokitty",
	}
	return animals[rand.Intn(len(animals)-1)]
}

func makeBasicHost(ctx context.Context, id peer.ID, addrs []ma.Multiaddr, ps pstore.Peerstore) (host.Host, error) {
	// Set up stream multiplexer
	tpt := msmux.NewBlankTransport()
	tpt.AddTransport("/yamux/1.0.0", yamux.DefaultTransport)

	// Create swarm (implements libP2P Network)
	swrm, err := swarm.NewSwarmWithProtector(
		ctx,
		addrs,
		id,
		ps,
		nil,
		tpt,
		nil,
	)
	if err != nil {
		return nil, err
	}

	netw := (*swarm.Network)(swrm)
	basicHost := bhost.New(netw)
	return basicHost, nil
}

// SendMessage opens a stream & sends a message from p to peerID
func (p *Peer) SendMessage(peerID peer.ID, msg Message) error {
	s, err := p.Host.NewStream(p.ctx, peerID, testbedProtocolID)
	if err != nil {
		return fmt.Errorf("error opening stream: %s", err.Error())
	}
	defer s.Close()

	ws := WrapStream(s)
	go p.handleStream(ws)
	return sendMessage(msg, ws)
}

// RandomPeer sends a message to a random peer
func (p Peer) RandomPeer() peer.ID {
	ps := p.Peerstore.Peers()
	peer := ps[rand.Intn(len(ps)-1)]
	if peer == p.ID {
		return p.RandomPeer()
	}
	return peer
}

// TestbedHandler is the handler we register with the multistream muxer
func (p Peer) TestbedHandler(s net.Stream) {
	defer s.Close()
	p.handleStream(WrapStream(s))
}

// handleStream is a for loop which receives and then sends a message.
// When Message.HangUp is true, it exits. This will close the stream
// on one of the sides. The other side's receiveMessage() will error
// with EOF, thus also breaking out from the loop.
func (p *Peer) handleStream(ws *WrappedStream) {
	for {
		// Read
		msg, err := receiveMessage(ws)
		if err != nil {
			if err.Error() == "EOF" {
				return
			}
			log.Errorf("error receiving message: %s", err.Error())
			return
		}

		log.Debugf("%s received message: %s", p.ID, msg.Type)
		switch msg.Type {
		case MtPing:
			if PingHandler(p, ws, msg) {
				return
			}
		case MtCLSnapshot:
			if ChandyLamportHandler(p, ws, msg) {
				return
			}
		}
	}
}
