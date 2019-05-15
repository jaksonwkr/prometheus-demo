package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/jaksonwkr/prometheus-demo/pkg/app/helloworld"
)

// HTTPHelloWorldHandler ...
type HTTPHelloWorldHandler struct {
	helloSvc          helloworld.HelloWorld
	durationHistogram *prometheus.HistogramVec
	totalRequests     *prometheus.CounterVec
}

// HelloWorldHandler handler HTTP RESTful for HelloWorld service
func (hdlr *HTTPHelloWorldHandler) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	var code int

	defer func(begun time.Time) {

		hdlr.durationHistogram.With(
			prometheus.Labels{"code": strconv.Itoa(code)},
		).Observe(time.Since(begun).Seconds())

		hdlr.totalRequests.With(prometheus.Labels{"code": strconv.Itoa(code)}).Inc()
	}(time.Now())

	resp, err := hdlr.helloSvc.SayHelloWorld()
	if err != nil {
		code = http.StatusInternalServerError
		w.WriteHeader(code)
		w.Write([]byte("error"))
		return
	}

	code = http.StatusOK
	w.WriteHeader(code)
	w.Write([]byte(resp))
	return
}

// NewHellowWorldHandler returns a new *HTTPHelloWorldHandler
func NewHellowWorldHandler(helloSvc helloworld.HelloWorld, durationHistogram *prometheus.HistogramVec,
	totalRequests *prometheus.CounterVec) *HTTPHelloWorldHandler {

	hdlr := &HTTPHelloWorldHandler{
		helloSvc:          helloSvc,
		durationHistogram: durationHistogram,
		totalRequests:     totalRequests,
	}

	return hdlr
}
