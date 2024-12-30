package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/okusarobert/toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type APIError struct {
	Code int
	Err  error
}

// Error implements the Error interface
func (e *APIError) Error() string {
	return e.Err.Error()
}

type HTTPFunc func(http.ResponseWriter, *http.Request) error

type HTTPMetricHandler struct {
	reqCounter prometheus.Counter
	errCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func NewHTTPMetricHandler(reqName string) *HTTPMetricHandler {
	reqCounter := promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
			Name:      "aggregator",
		},
	)
	errCounter := promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: fmt.Sprintf("http_%s_%s", reqName, "err_counter"),
			Name:      "aggregator",
		},
	)
	reqLatency := promauto.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_latency"),
			Name:      "aggregator",
			Buckets:   []float64{0.1, 0.5, 1},
		},
	)
	return &HTTPMetricHandler{
		reqCounter: reqCounter,
		reqLatency: reqLatency,
		errCounter: errCounter,
	}
}

func makeHTTPHandlerFunc(fn HTTPFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			if apiErr, ok := err.(*APIError); ok {
				writeJSON(w, apiErr.Code, map[string]string{"error": apiErr.Error()})
			}

		}
	}
}

func (h *HTTPMetricHandler) instrument(next HTTPFunc) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var err error
		defer func(start time.Time) {
			latency := time.Since(start).Seconds()
			logrus.WithFields(logrus.Fields{
				"latency": latency,
				"request": r.RequestURI,
				"error":   err,
			}).Info("http requests")
			h.reqLatency.Observe(latency)
			if err != nil {
				h.errCounter.Inc()
			}
		}(time.Now())
		h.reqCounter.Inc()
		err = next(w, r)
		return err
	}
}

func handleGetInvoice(svc Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != "GET" {
			return &APIError{
				Code: http.StatusMethodNotAllowed,
				Err:  fmt.Errorf("method %s not allowed", r.Method),
			}
		}
		obuID := r.URL.Query().Get("obu")
		if obuID == "" {
			return &APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid OBU ID %s", obuID),
			}
		} else {
			id, err := strconv.Atoi(obuID)
			if err != nil {
				return &APIError{
					Code: http.StatusInternalServerError,
					Err:  fmt.Errorf("internal server error"),
				}
			}
			invoice, err := svc.CalculateInvoice(id)
			if err != nil {
				return &APIError{
					Code: http.StatusInternalServerError,
					Err:  err,
				}
			}
			return writeJSON(w, http.StatusOK, invoice)
		}
	}
}

func handleAggregate(svc Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var distance types.Distance
		if r.Method != "POST" {
			return &APIError{
				Code: http.StatusMethodNotAllowed,
				Err:  fmt.Errorf("method %s not supported", r.Method),
			}
		}
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			return &APIError{
				Code: http.StatusMethodNotAllowed,
				Err:  err,
			}
		}
		if err := svc.AggregateDistance(distance); err != nil {
			return &APIError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}
		return nil
	}
}
