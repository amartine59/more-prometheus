package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	totalRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get request",
		},
		[]string{"path"},
	)

	responseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status_total",
			Help: "Status of HTTP response",
		},
		[]string{"status"},
	)

	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_response_time_seconds",
			Help: "Duration of HTTP Requests",
			// Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func registerCustomMetrics() error {
	err := prometheus.Register(totalRequestCount)
	if err != nil {
		return fmt.Errorf("failed to register %s metric: %w", "Total Request Count", err)
	}

	err = prometheus.Register(responseStatus)
	if err != nil {
		return fmt.Errorf("failed to register %s metric: %w", "Response Status Count", err)
	}

	err = prometheus.Register(httpDuration)
	if err != nil {
		return fmt.Errorf("failed to register %s metric: %w", "Http duration histogram", err)
	}

	return nil
}

func prometheusMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		next.ServeHTTP(w, r)

		fmt.Println("Entire Map: ", w.Header())

		// Still not sure how to catch the response status without creating my own response writer
		responseStatus.WithLabelValues(strconv.Itoa(http.StatusOK)).Inc()
		totalRequestCount.WithLabelValues(path).Inc()
		timer.ObserveDuration()
	})
}

func main() {
	err := registerCustomMetrics()
	if err != nil {
		// this only makes sense for our exercise
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.Use(prometheusMW)

	// Reporting metrics for Prometheus using default handler
	r.Path("/metrics").Handler(promhttp.Handler())

	r.PathPrefix("/").
		Handler(
			http.FileServer(http.Dir("./static/")),
		)

	fmt.Println("Starting App on port 2112")
	err = http.ListenAndServe(":2112", r)
	if err != nil {
		// Should kill it if something goes wrong
		log.Fatal(err)
	}
}
