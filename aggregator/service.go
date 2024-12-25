package main

import (
	"fmt"

	"github.com/okusarobert/toll-calculator/types"
)

const (
	basePrice = 3.15
)

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storage interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storage
}

func NewInvoiceAggregator(store Storage) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Println("processing and inserting distance into storage: ", distance)
	return i.store.Insert(distance)
}

func (i *InvoiceAggregator) CalculateInvoice(id int) (*types.Invoice, error) {
	dist, err := i.store.Get(id)
	if err != nil {
		return nil, err
	}
	return &types.Invoice{
		OBUID:         id,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}, nil
}
