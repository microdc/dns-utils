package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func usage() {
	fmt.Println("dns-utils -h for usage.")
	os.Exit(1)
}

var (
	lookupLatencyHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "dns_lookup_latency",
		Help:    "Latency of DNS lookups in seconds",
		Buckets: []float64{1, 2, 5, 10},
	})

	failedLookupCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dns_failed_lookups",
		Help: "The total number of failed DNS lookups",
	})
)

const timeout = 5 * time.Second

func lookup(domain string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	t0 := time.Now()
	var r net.Resolver
	_, err := r.LookupHost(ctx, domain)
	if err != nil {
		fmt.Println("TIMEOUT")
		failedLookupCounter.Inc()
	}
	t1 := time.Now()
	duration := t1.Sub(t0)

	lookupLatencyHistogram.Observe(duration.Seconds())

	fmt.Println(duration)
}

var domain string
var interval int

func main() {
	flag.StringVar(&domain, "domain", "", "domain to resolve")
	flag.IntVar(&interval, "interval", 10, "number of concurrent requests to send")
	flag.Parse()

	if domain == "" {
		usage()
	}

	go func(domain string) {
		for i := 1; ; i++ {

			go lookup(domain)
			if i%interval == 0 {
				time.Sleep(time.Second * 5)
			}
		}
	}(domain)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
