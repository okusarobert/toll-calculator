package main

import (
	"time"

	"github.com/okusarobert/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

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
