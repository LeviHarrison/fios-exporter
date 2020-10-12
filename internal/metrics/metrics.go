package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// TXMinute1 is the average kilobytes of traffic per second transmitted in the last minute
	TXMinute1 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "fios_tx_minute_1",
		Help: "Average kilobytes of traffic per second transmitted in the last minute.",
	})

	// RXMinute1 is the average kilobytes of traffic per second recieved in the last minute
	RXMinute1 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "fios_rx_minute_1",
		Help: "Average kilobytes of traffic per second recieved in the last minute.",
	})
)
