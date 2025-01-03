package main

import (
	"log"

	"github.com/okusarobert/toll-calculator/aggregator/client"
)

const kafkaTopic = "obudata"

// Transport (HTTP, GRPC, Kafka) -> attach business logic

func main() {
	svc := NewCalculatorService()
	svc = NewLogMiddleware(svc)
	httpClient := client.NewHTTPClient("http://localhost:4000")
	grpcClient, err := client.NewGRPCClient(":3001")
	if err != nil {
		log.Fatal(err)
	}
	_ = grpcClient
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, httpClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
