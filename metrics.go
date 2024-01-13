package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	httpClient *http.Client

	// It's for local prometheus only
	scrape_failures *prometheus.CounterVec

	mpc_state prometheus.Gauge

	timestamps     *prometheus.CounterVec
	lastTimestamps map[string]float64

	heights     *prometheus.CounterVec
	lastHeights map[string]float64
	mutex       sync.Mutex
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		httpClient: &http.Client{},
		timestamps: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "metis:sequencer:timestamp",
				Help: "Current Unix timestamp of the service.",
			},
			[]string{"svc_name"},
		),
		heights: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "metis:sequencer:height",
				Help: "Current block number of the service.",
			},
			[]string{"svc_name"},
		),
		mpc_state: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "metis:sequencer:mpc:state",
			Help: "the mpc service signature service is working or not.",
		}),
		scrape_failures: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "metis_sequencer_exporter_failures",
				Help: "Number of scrape errors.",
			},
			[]string{"url"},
		),
		lastTimestamps: make(map[string]float64),
		lastHeights:    make(map[string]float64),
	}
	reg.MustRegister(m.timestamps)
	reg.MustRegister(m.heights)
	reg.MustRegister(m.scrape_failures)
	reg.MustRegister(m.mpc_state)
	return m
}

func (m *Metrics) ScrapeSequencerState(basectx context.Context, url string) {
	ticker := time.NewTimer(0)
	defer ticker.Stop()

	scrape := func() error {
		newctx, cancel := context.WithTimeout(basectx, time.Second*3)
		defer cancel()

		req, err := http.NewRequestWithContext(newctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		resp, err := m.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var state NodeHealthResp
		if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
			return err
		}

		log.Println("seq:height", state.L2.BlockNumber, "seq:timetamp", state.L2.Timestamp)

		m.mutex.Lock()
		defer m.mutex.Unlock()

		if t := state.L2.Timestamp - m.lastTimestamps["seq"]; t > 0 {
			m.timestamps.With(prometheus.Labels{"svc_name": "seq"}).Add(t)
			m.lastTimestamps["seq"] += t
		}

		if t := state.L2.BlockNumber - m.lastHeights["seq"]; t > 0 {
			m.heights.With(prometheus.Labels{"svc_name": "seq"}).Add(t)
			m.lastHeights["seq"] += t
		}

		if state.MPC != nil && state.MPC.IsMpcProposer == 1 {
			log.Println("mpc:timestmap", state.MPC.Timestamp, "mpc:signSuccess", state.MPC.SignSuccess)

			if t := state.MPC.Timestamp - m.lastTimestamps["mpc"]; t > 0 {
				m.timestamps.With(prometheus.Labels{"svc_name": "mpc"}).Add(t)
				m.lastTimestamps["mpc"] += t
			}

			m.mpc_state.Set(float64(state.MPC.SignSuccess))
		}

		return nil
	}

	for {
		select {
		case <-basectx.Done():
			return
		case <-ticker.C:
			if err := scrape(); err != nil {
				m.scrape_failures.With(prometheus.Labels{"url": url}).Inc()
				log.Println("Failed to scrape sequencer state", err)
			}
			ticker.Reset(time.Minute)
		}
	}
}

func (m *Metrics) ScrapeNodeState(basectx context.Context, url string) {
	ticker := time.NewTimer(0)
	defer ticker.Stop()

	scrape := func() error {
		newctx, cancel := context.WithTimeout(basectx, time.Second*3)
		defer cancel()

		req, err := http.NewRequestWithContext(newctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		resp, err := m.httpClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		var state LastestSpanResp
		if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
			return err
		}

		log.Println("node:height", state.Height)
		m.mutex.Lock()
		defer m.mutex.Unlock()

		if t := state.Height - m.lastHeights["node"]; t > 0 {
			m.heights.With(prometheus.Labels{"svc_name": "node"}).Add(t)
			m.lastHeights["node"] += t
		}
		return nil
	}

	for {
		select {
		case <-basectx.Done():
			return
		case <-ticker.C:
			if err := scrape(); err != nil {
				m.scrape_failures.With(prometheus.Labels{"url": url}).Inc()
				log.Println("Failed to scrape the node state", err)
			}
			ticker.Reset(time.Minute)
		}
	}
}
