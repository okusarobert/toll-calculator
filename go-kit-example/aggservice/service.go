package aggservice

import (
	"context"
	"fmt"

	"github.com/go-kit/log"
	"github.com/okusarobert/toll-calculator/types"
)

const (
	basePrice = 35.99
)

type Service interface {
	Aggregate(context.Context, types.Distance) error
	Calculate(context.Context, int) (*types.Invoice, error)
}

type BasicService struct {
	store Storage
}

func newBasicService(store Storage) Service {
	return &BasicService{
		store: store,
	}
}

func (s *BasicService) Aggregate(ctx context.Context, distance types.Distance) error {
	fmt.Println("this is coming from the internal businness logic layer")
	// logrus.WithFields(logrus.Fields{
	// 	"obuID": distance.OBUID,
	// 	"unix":  distance.Unix,
	// 	"value": distance.Value,
	// }).Info("Aggregate")
	return s.store.Insert(distance)
}

func (s *BasicService) Calculate(ctx context.Context, id int) (*types.Invoice, error) {
	dist, err := s.store.Get(id)
	if err != nil {
		return nil, err
	}
	return &types.Invoice{
		OBUID:         id,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}, nil
}

// NewAggregatorService will construct a complete microservice with logging and instrumentation middleware
func New(logger log.Logger) Service {
	var svc Service
	{
		svc = newBasicService(NewMemoryStore())
		svc = newLoggingMiddleware(logger)(svc)
		svc = newInstrumentationMiddleware(logger)(svc)
	}
	return svc
}
