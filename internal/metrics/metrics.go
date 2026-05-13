package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	URLsShortenedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "snipqurl_urls_shortened_total",
			Help: "Total number of URLs shortened",
		},
	)

	QRGeneratedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "snipqurl_qr_generated_total",
			Help: "Total number of QR codes generated",
		},
	)

	RedirectsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "snipqurl_redirects_total",
			Help: "Total number of redirects",
		},
	)
)
