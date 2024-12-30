package main

import (
	"fmt"
	"time"

	"github.com/okusarobert/toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type MetricsMiddleware struct {
	reqCounterAgg  prometheus.Counter
	reqCounterCalc prometheus.Counter
	errCounterAgg  prometheus.Counter
	errCounterCalc prometheus.Counter
	reqLatencyAgg  prometheus.Histogram
	reqLatencyCalc prometheus.Histogram
	next           Aggregator
}

func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter_agg",
		Name:      "aggregate",
	})
	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter_calc",
		Name:      "calc",
	})
	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "calc",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	errCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_err_counter_agg",
		Name:      "aggregate",
	})
	errCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_err_counter_calc",
		Name:      "aggregate",
	})
	return &MetricsMiddleware{
		next:           next,
		reqCounterAgg:  reqCounterAgg,
		reqCounterCalc: reqCounterCalc,
		reqLatencyCalc: reqLatencyCalc,
		reqLatencyAgg:  reqLatencyAgg,
		errCounterAgg:  errCounterAgg,
		errCounterCalc: errCounterCalc,
	}
}

func (m *MetricsMiddleware) AggregateDistance(distance types.Distance) (err error) {

	defer func(start time.Time) {
		m.reqLatencyAgg.Observe(time.Since(start).Seconds())
		m.reqCounterAgg.Inc()
		if err != nil {
			m.errCounterAgg.Inc()
		}
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return

}
func (m *MetricsMiddleware) CalculateInvoice(id int) (invoice *types.Invoice, err error) {

	defer func(start time.Time) {
		fmt.Println("calculating invoice ...")
		m.reqLatencyCalc.Observe(time.Since(start).Seconds())
		m.reqCounterCalc.Inc()
		if err != nil {
			m.errCounterCalc.Inc()
		}
	}(time.Now())
	invoice, err = m.next.CalculateInvoice(id)
	return
}

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (i *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":  time.Since(start),
			"error": err,
			"value": distance.Value,
			"obuID": distance.OBUID,
			"unix":  distance.Unix,
		}).Info("aggregate")
	}(time.Now())
	err = i.next.AggregateDistance(distance)
	return

}

func (i *LogMiddleware) CalculateInvoice(id int) (invoice *types.Invoice, err error) {
	var (
		amount   float64
		distance float64
	)
	if invoice != nil {
		amount = invoice.TotalAmount
		distance = invoice.TotalDistance
	}
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"error":    err,
			"obuID":    id,
			"amount":   amount,
			"distance": distance,
		}).Info("CalculateInvoice")
	}(time.Now())
	invoice, err = i.next.CalculateInvoice(id)
	return

}
