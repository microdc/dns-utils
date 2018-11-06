package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func usage() {
	fmt.Println("dns-utils -h for usage.")
	os.Exit(1)
}

func lookup(domain string) {
	t0 := time.Now()
	_, err := net.LookupHost(domain)
	if err != nil {
		panic("error")
	}
	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
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

	go func() {
		for i := 1; ; i++ {
			go lookup(domain)
			if i%interval == 0 {
				time.Sleep(time.Second * 5)
			}
		}
	}()

	select {}
}
