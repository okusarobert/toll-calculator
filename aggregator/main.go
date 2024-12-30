package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/okusarobert/toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	var (
		store          = makeStore()
		svc            = NewInvoiceAggregator(store)
		grpcListenAddr = os.Getenv("AGG_GRPC_ENDPOINT")
		httpListenAddr = os.Getenv("AGG_HTTP_ENDPOINT")
	)
	svc = NewMetricsMiddleware(svc)
	svc = NewLogMiddleware(svc)
	go func() {
		log.Fatal(makeGRPCTransport(grpcListenAddr, svc))
	}()
	log.Fatal(makeHTTPTransport(httpListenAddr, svc))
}

func makeStore() Storage {
	store := os.Getenv("AGG_STORE_TYPE")
	switch store {
	case "memory":
		return NewMemoryStore()
	default:
		log.Fatalf("invalid store type given %s", store)
		return nil
	}
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("GRPC transport running on port", listenAddr)
	// make a TCP listener
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	// Make a new GRPC native server with (options)
	server := grpc.NewServer([]grpc.ServerOption{}...)
	//  Register our GRPC implementation to the GRPC package
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	aggregateMetricHandler := NewHTTPMetricHandler("aggregate")
	invoiceMetricHandler := NewHTTPMetricHandler("invoice")
	http.HandleFunc("/aggregate", makeHTTPHandlerFunc(aggregateMetricHandler.instrument(handleAggregate(svc))))
	http.HandleFunc("/invoice", makeHTTPHandlerFunc(invoiceMetricHandler.instrument(handleGetInvoice(svc))))
	// http.HandleFunc("/invoice", invoiceMetricHandler.instrument(handleGetInvoice(svc)))
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("HTTP transport running on port", listenAddr)
	return http.ListenAndServe(listenAddr, nil)
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
