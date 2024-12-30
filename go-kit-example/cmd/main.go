package main

import (
	"net"
	"net/http"
	"os"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/log"
	"github.com/okusarobert/toll-calculator/go-kit-example/aggendpoint"
	"github.com/okusarobert/toll-calculator/go-kit-example/aggservice"
	"github.com/okusarobert/toll-calculator/go-kit-example/aggtransport"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	httpListener, err := net.Listen("tcp", ":3004")
	if err != nil {
		logger.Log("transport", "HTTP", "during", "Listen", "err", err)
		os.Exit(1)
	}
	var duration metrics.Histogram
	{
		// Endpoint-level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "aggregator",
			Subsystem: "aggservice",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}
	var (
		service     = aggservice.New(logger)
		endpoints   = aggendpoint.New(service, logger, duration)
		httpHandler = aggtransport.NewHTTPHandler(endpoints, logger)
	)
	logger.Log("transport", "HTTP", "addr", ":3004")
	err = http.Serve(httpListener, httpHandler)
	if err != nil {
		panic(err)
	}
}
