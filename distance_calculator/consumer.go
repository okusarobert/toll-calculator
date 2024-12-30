package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/okusarobert/toll-calculator/aggregator/client"
	"github.com/okusarobert/toll-calculator/types"
	"github.com/sirupsen/logrus"
)

// This can also be called kafka transport
type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
	aggClient   client.Client
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, aggClient client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}
	err = c.SubscribeTopics([]string{topic}, nil)

	if err != nil {
		return nil, err
	}
	return &KafkaConsumer{
		consumer:    c,
		calcService: svc,
		aggClient:   aggClient,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("Kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {
	defer c.consumer.Close()
	// A signal handler or similar could be used to set this to false to break the loop.
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error %s", err.Error())
			continue
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialization error: %s", err.Error())
			logrus.WithFields(logrus.Fields{
				"error":     err,
				"requestID": data.RequestID,
			})
			continue
		}
		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("error calculating distance: %s", err.Error())
			continue
		}
		req := &types.AggregateRequest{
			Value: distance,
			ObuID: int32(data.OBUID),
			Unix:  time.Now().UnixNano(),
		}
		if err = c.aggClient.Aggregate(context.Background(), req); err != nil {
			logrus.Error(err)
			continue
		}
	}
}
