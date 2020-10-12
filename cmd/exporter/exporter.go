package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/leviharrison/fios-exporter/internal/metrics"
	"github.com/leviharrison/fios-exporter/internal/scrape"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	host     = flag.String("host", "https://myfiosgateway.com", "The address of your router")
	password = flag.String("password", "", "The password of your router")
	port     = flag.Int("port", 2190, "The port which the exporter listens on")
)

func init() {
	flag.Parse()

	if *password == "" {
		log.Fatal("You must have a password")
	}

	scrape.Init(*host)

	prometheus.MustRegister(metrics.TXMinute1)
	prometheus.MustRegister(metrics.RXMinute1)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	go scrape.Scrape(*host, *password)

	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(*port), nil))
}
