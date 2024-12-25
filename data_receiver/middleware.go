package main

import (
	"fmt"
	"time"

	"github.com/okusarobert/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) DataProducer {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) ProduceData(data types.OBUData) error {
	fmt.Println("logging middleware ...")
	start := time.Now()
	defer func(start time.Time) {
		logrus.WithFields(
			logrus.Fields{
				"obuID": data.OBUID,
				"lat":   data.Lat,
				"long":  data.Long,
				"took":  time.Since(start),
			},
		).Info("producing to kafka")
	}(start)
	return l.next.ProduceData(data)
}
