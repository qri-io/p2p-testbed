package main

import (
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
)

const (
	// MtCLSnapshot represents a chandy-lamport snapshot message
	MtCLSnapshot = MsgType("cl_snapshot")
	clSnapshots  = "cl_snapshots"
)

// InitiateCLSnapshot kicks off a chandy-lamport snapshot
// This Peer that's initiating the process:
//     Saves its own local state
//     Sends a snapshot request message bearing a snapshot token to all other processes
// A process receiving the snapshot token for the first time on any message:
//     Sends the observer process its own saved state
//     Attaches the snapshot token to all subsequent messages (to help propagate the snapshot token)
// When a process that has already received the snapshot token receives a message that does not bear the snapshot token,
// this process will forward that message to the observer process.
// This message was obviously sent before the snapshot “cut off” (as it does not bear a snapshot token and thus must
// have come from before the snapshot token was sent out) and needs to be included in the snapshot.
func InitiateCLSnapshot(p *Peer) {
	log.Infof("%s initiating CL snapshot", p.ID)
	snapshotToken := NewMessageID()

	span := opentracing.StartSpan("cl_snapshot")
	p.SetState(clSpanStateKey(snapshotToken), span)

	TakeCLSnapshot(p, snapshotToken, span)
}

// TakeCLSnapshot checks if a snapshot has been taken,
// if not it propagates the snapshot token to all connected peers
// and takes a local snapshot
func TakeCLSnapshot(p *Peer, snapshotToken string, span opentracing.Span) {
	// get snapshots map from state
	snapshots := map[string]interface{}{}
	if val := p.GetState(clSnapshots); val != nil {
		if sl, ok := val.(map[string]interface{}); ok {
			snapshots = sl
		}
	}

	// check if snapshot is empty
	if snapshots[snapshotToken] == nil {
		log.Infof("%s recording snapshot. forwarding to %d peers", p.ID, len(p.Peerstore.Peers())-1)
		// record snapshot
		ss := LocalCLSnapshot(p)
		snapshots[snapshotToken] = ss
		p.SetState(clSnapshots, snapshots)
		span.SetTag("state", ss)

		// send marker to each connected peer
		for _, peer := range p.Peerstore.Peers() {
			if peer != p.ID {
				p.SendMessage(peer, NewMessage(MtCLSnapshot, snapshotToken, span))
			}
		}
	}
}

// LocalCLSnapshot generates a snapshot of a peer's local state
func LocalCLSnapshot(p *Peer) string {
	conns := p.Host.Network().Conns()
	connPeers := make([]string, len(conns))
	for i, c := range conns {
		connPeers[i] = c.RemotePeer().Pretty()
	}

	state := map[string]interface{}{
		"peerID": p.ID.Pretty(),
		"pet":    p.GetState("pet"),
		"conns":  connPeers,
	}

	data, err := json.Marshal(state)
	if err != nil {
		log.Errorf("error marshaling state data to json: %s", err.Error())
	}

	return string(data)
}

// ChandyLamportHandler handles messages of type MtCLSnapshot
func ChandyLamportHandler(p *Peer, ws *WrappedStream, msg Message) (hangup bool) {
	if snapshotToken, ok := msg.Payload.(string); ok {
		var span opentracing.Span
		wireContext, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, msg.Tracing)
		if err != nil {
			// If for whatever reason a span isn't found, go ahead an start a new root span
			span = opentracing.StartSpan("snapshot")
		} else {
			span = opentracing.StartSpan("snapshot", opentracing.ChildOf(wireContext))
		}
		TakeCLSnapshot(p, snapshotToken, span)
		span.Finish()

		// if we initiated this span, close it off
		if initspan, ok := p.GetState(clSpanStateKey(snapshotToken)).(opentracing.Span); ok {
			// log.Infof("%s finalized snapshot: %s", p.ID, snapshotToken)
			initspan.Finish()
		}
	}
	return true
}

// clSpanStateKey is where we keep spans in state for initiated snapshots
func clSpanStateKey(snapshotToken string) string {
	return fmt.Sprintf("initiatedSnapshot.%s", snapshotToken)
}
