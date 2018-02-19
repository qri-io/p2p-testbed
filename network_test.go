package main

import (
	"context"
	"testing"
)

func TestNetwork(t *testing.T) {
	ctx := context.Background()
	peers, err := NewNetwork(ctx, 2)
	if err != nil {
		t.Errorf("error creating network: %s", err.Error())
		return
	}

	if err := connectPeers(ctx, peers); err != nil {
		t.Errorf("error connecting peers: %s", err.Error())
	}
}
