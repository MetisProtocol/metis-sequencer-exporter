package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		SeqStateURL  string
		NodeStateURL string
	)

	flag.StringVar(&SeqStateURL, "url.state.seq", "http://localhost:9545/health", "the sequencer state url")
	flag.StringVar(&NodeStateURL, "url.state.node", "http://localhost:1317/metis/latest-span", "the node state url")
	flag.Parse()

	basectx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	reg := prometheus.NewRegistry()

	metrics := NewMetrics(reg)

	go metrics.ScrapeSequencerState(basectx, SeqStateURL)
	go metrics.ScrapeNodeState(basectx, NodeStateURL)

	server := &http.Server{Addr: ":21012"}
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	go func() {
		<-basectx.Done()
		_ = server.Shutdown(context.Background())
	}()

	log.Println("serving")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
