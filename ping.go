package main

import (
	"math/rand"
	"time"

	"github.com/opentracing/opentracing-go"

	peer "gx/ipfs/QmXYjuNuxVzXKJCfWasQk1RqkhVLDM9jtUKhqc2WPQmFSB/go-libp2p-peer"
)

const (
	// MtPing is a ping/pong message
	MtPing = MsgType("PING")
)

// SendPing initiates a ping message from peer to a known peer.ID
func SendPing(p *Peer, peerID peer.ID) error {
	span := opentracing.StartSpan("PING")
	msg := NewMessage(MtPing, "PING", span)
	// add our span to state according to message ID so we can grab it
	p.SetState(msg.ID, span)
	// simulate network latency by sleeping for a random number of milliseconds
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	return p.SendMessage(peerID, msg)
}

// PingHandler handles messages of type MtPing
func PingHandler(p *Peer, ws *WrappedStream, msg Message) (hangup bool) {
	if payload, ok := msg.Payload.(string); ok {
		switch payload {
		case "PONG":
			// grab the span out of state & close it
			if span, ok := p.GetState(msg.ID).(opentracing.Span); ok {
				span.Finish()
			}
		case "PING":
			var span opentracing.Span
			wireContext, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, msg.Tracing)
			if err != nil {
				// If for whatever reason a span isn't found, go ahead an start a new root span
				span = opentracing.StartSpan("PONG")
			} else {
				span = opentracing.StartSpan("PONG", opentracing.ChildOf(wireContext))
			}
			defer span.Finish()

			// simulate network latency by sleeping for a random number of milliseconds
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			p.SendMessage(ws.stream.Conn().RemotePeer(), Message{
				ID:      msg.ID,
				Tracing: msg.Tracing,
				Type:    MtPing,
				Payload: "PONG",
			})
		}
	}
	return true
}
