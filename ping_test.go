package main

import (
	"context"
	"testing"
)

func TestPing(t *testing.T) {
	ctx := context.Background()
	peers, err := NewNetwork(ctx, 50)
	if err != nil {
		t.Errorf("error creating network: %s", err.Error())
		return
	}

	if err := connectPeers(ctx, peers); err != nil {
		t.Errorf("error connecting peers: %s", err.Error())
	}

	for i, p1 := range peers {
		for _, p2 := range peers[i+1:] {
			go func() {
				if err := SendPing(p1, p2.ID); err != nil {
					t.Errorf("%s -> %s error: %s", p1.ID.Pretty(), p2.ID.Pretty(), err.Error())
					return
				}
			}()
		}
	}
}
