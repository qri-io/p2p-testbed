package main

import (
	"context"
	"testing"
	"time"
)

func TestCLSnapshot(t *testing.T) {
	ctx := context.Background()
	done := make(chan bool, 1)
	timer := time.NewTimer(time.Second * 10)
	peers, err := NewNetwork(ctx, 10)
	callbackCalled := false

	if err != nil {
		t.Errorf("error creating network: %s", err.Error())
		return
	}

	if err := connectPeers(ctx, peers); err != nil {
		t.Errorf("error connecting peers: %s", err.Error())
	}

	go func() {
		<-timer.C
		t.Errorf("snapshot took too long")
		done <- true
	}()

	callback := func(snapshotToken string) {
		if callbackCalled {
			t.Errorf("callback called multiple times")
			return
		}
		callbackCalled = true

		missed := 0
		for _, p := range peers {
			// get snapshots map from state
			snapshots := map[string]interface{}{}
			if val, ok := p.state.Load(clSnapshots); ok {
				if sl, ok := val.(map[string]interface{}); ok {
					snapshots = sl
				}
			}

			// check if snapshot is empty
			if snapshots[snapshotToken] == nil {
				t.Logf("peer %s didn't record snapshot", p.ID)
				missed++
			}
		}
		if missed > 0 {
			t.Errorf("snapshot missed %d/%d peers", missed, len(peers))
		}
		done <- true
	}

	InitiateCLSnapshot(peers[0], callback)
	<-done
}
