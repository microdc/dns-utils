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
		Name:    "dns_lookup_duration_seconds",
		Help:    "Duration of DNS lookups in seconds",
		Buckets: []float64{0.005, 0.01, 0.05, 0.1, 0.15, 0.20, 0.25, 0.3, 0.35, 0.4, 0.45, 0.5, 1, 2},
	})

	failedLookupCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dns_lookup_error",
		Help: "The total number of failed DNS lookups",
	})
)

const timeout = 5 * time.Second

func lookup(domain string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	var r net.Resolver
	_, err := r.LookupHost(ctx, domain)
	if err != nil {
		failedLookupCounter.Inc()
	}
	end := time.Now()
	duration := end.Sub(start)

	lookupLatencyHistogram.Observe(duration.Seconds())

	fmt.Println(duration.Seconds())
}

func sampleLookups(domain string, interval int) {
	go func() {
		for i := 1; ; i++ {
			go lookup(domain)
			if i%interval == 0 {
				time.Sleep(time.Second * 5)
			}
		}
	}()
}

var domain string
var interval int
var port int

func main() {
	flag.StringVar(&domain, "domain", "", "domain to resolve")
	flag.IntVar(&interval, "interval", 10, "number of concurrent requests to send")
	flag.IntVar(&port, "port", 8080, "port number to listen on (for prometheus scraping)")
	flag.Parse()

	if domain == "" {
		usage()
	}

	sampleLookups(domain, interval)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
