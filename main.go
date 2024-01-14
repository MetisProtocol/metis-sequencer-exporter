package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		SeqStateURL    string
		NodeStateURL   string
		ScrapeInterval time.Duration
	)

	flag.StringVar(&SeqStateURL, "url.state.seq", "http://localhost:9545/health", "the sequencer state url")
	flag.StringVar(&NodeStateURL, "url.state.node", "http://localhost:1317/metis/latest-span", "the node state url")
	flag.DurationVar(&ScrapeInterval, "scrape.interval", time.Second*15, "scrape interval")
	flag.Parse()

	basectx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	reg := prometheus.NewRegistry()

	metrics := NewMetrics(reg)

	go metrics.ScrapeSequencerState(basectx, SeqStateURL, ScrapeInterval)
	go metrics.ScrapeNodeState(basectx, NodeStateURL, ScrapeInterval)

	server := &http.Server{Addr: ":21012"}
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "pong") })

	go func() {
		log.Println("serving")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cancel()
			log.Println(err)
		}
	}()

	<-basectx.Done()
	log.Println("graceful stopping")
	_ = server.Shutdown(context.Background())
}
