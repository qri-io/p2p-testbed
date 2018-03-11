package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"

	logging "github.com/ipfs/go-log"
	// logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"

	opentracing "github.com/opentracing/opentracing-go"
	"sourcegraph.com/sourcegraph/appdash"
	appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"
	"sourcegraph.com/sourcegraph/appdash/traceapp"
)

var (
	peerCount = 20
	log       = logging.Logger("testbed")
)

func init() {
	logging.Configure(logging.Output(os.Stdout), logging.LevelInfo)
}

func main() {
	startAppDash()

	log.Infof("starting network with %d peers", peerCount)
	ctx := context.Background()
	peers, err := Setup(ctx, peerCount)
	if err != nil {
		log.Fatalf("error setting up network: %s", err.Error())
	}

	// send a ping
	if err := SendPing(peers[0], peers[0].RandomPeer()); err != nil {
		log.Errorf("error sending message: %s", err.Error())
	}

	// take a CL snapshot
	InitiateCLSnapshot(peers[1], func(snapshotToken string) {
		log.Infof("recorded snapshot: %s", snapshotToken)
	})

	// block forever
	<-make(chan bool)
}

func startAppDash() {
	store := appdash.NewMemoryStore()
	// Listen on any available TCP port locally.
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		log.Fatal(err)
	}
	collectorPort := l.Addr().(*net.TCPAddr).Port
	collectorAdd := fmt.Sprintf(":%d", collectorPort)

	// Start an Appdash collection server that will listen for spans and
	// annotations and add them to the local collector (stored in-memory).
	cs := appdash.NewServer(l, appdash.NewLocalCollector(store))
	go cs.Start()

	// Print the URL at which the web UI will be running.
	appdashPort := 8700
	appdashURLStr := fmt.Sprintf("http://localhost:%d", appdashPort)
	appdashURL, err := url.Parse(appdashURLStr)
	if err != nil {
		log.Fatalf("Error parsing %s: %s", appdashURLStr, err)
	}
	log.Infof("To see your traces, go to %s/traces", appdashURL)

	// Start the web UI in a separate goroutine.
	tapp, err := traceapp.New(nil, appdashURL)
	if err != nil {
		log.Fatal(err)
	}
	tapp.Store = store
	tapp.Queryer = store
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appdashPort), tapp))
	}()

	tracer := appdashot.NewTracer(appdash.NewRemoteCollector(collectorAdd))
	opentracing.InitGlobalTracer(tracer)
}
