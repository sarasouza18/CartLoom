package utils

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// InitMetricsServer initializes the metrics endpoint for Prometheus
func InitMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("Serving metrics on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Error serving metrics: %v", err)
		}
	}()
}
