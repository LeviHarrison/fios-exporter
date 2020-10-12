package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// TXMinute1 is the average kilobits of traffic per second transmitted in the last minute
	TXMinute1 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "fios_tx_kbs_minute_1",
		Help: "Average kilobits of traffic per second transmitted in the last minute.",
	})

	// RXMinute1 is the average kilobits of traffic per second received in the last minute
	RXMinute1 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "fios_rx_kbs_minute_1",
		Help: "Average kilobits of traffic per second received in the last minute.",
	})
)
