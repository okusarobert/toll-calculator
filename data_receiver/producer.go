package main

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/okusarobert/toll-calculator/types"
)

type DataProducer interface {
	ProduceData(types.OBUData) error
}

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProducer(topic string) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		return nil, err
	}
	// defer p.Close()
	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			fmt.Println(e)
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return &KafkaProducer{
		producer: p,
		topic:    topic,
	}, nil
}

func (p *KafkaProducer) ProduceData(data types.OBUData) error {
	// defer p.producer.Close()
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Println("Producing data : ", p.producer.IsClosed())
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny},
		Value: b,
	}, nil)
	fmt.Printf("Error: %+v\n", err)
	p.producer.Flush(15 * 1000)
	return err
}
