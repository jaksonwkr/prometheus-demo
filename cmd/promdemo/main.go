package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jaksonwkr/prometheus-demo/pkg/app/handlers"
	"github.com/jaksonwkr/prometheus-demo/pkg/app/helloworld"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	// Initialize Prometheus Metrics
	durationHistogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Help: "Histogram for the requests duration",
			Name: "http_requests_duration_buckets",
			// Buckets: []float64{.005, .01, .025, .05, .1, .25},
		},
		[]string{"code"},
	)

	totalRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Help: "Total of requests of the server",
			Name: "http_requests_total",
		},
		[]string{"code"},
	)

	promRegister := prometheus.NewRegistry()
	promRegister.MustRegister(durationHistogram, totalRequests)

	// Instantiate HelloWorld service
	helloSvc, err := helloworld.New()
	if err != nil {
		fmt.Printf("Could not instantiate HelloWorld service: %s\n", err)
		return
	}

	// Instantiate the HelloWorldHandler that implements the http.Handler
	helloWorldHandler := handlers.NewHellowWorldHandler(*helloSvc, durationHistogram, totalRequests)

	// Add the routes to the Handler
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(promRegister, promhttp.HandlerOpts{}))
	mux.HandleFunc("/helloworld", helloWorldHandler.HelloWorldHandler)

	// Starts HTTP Server
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
		Addr:         ":8080",
		Handler:      mux,
	}

	fmt.Println("Listening on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
