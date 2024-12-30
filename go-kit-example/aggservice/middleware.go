package aggservice

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/okusarobert/toll-calculator/types"
)

type Middleware func(Service) Service

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func newLoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

func (mw *loggingMiddleware) Aggregate(ctx context.Context, distance types.Distance) (err error) {
	defer func(start time.Time) {
		mw.logger.Log("method", "Aggregate", "error", err, "took", time.Since(start),
			"obuID", distance.OBUID, "value", distance.Value, "unix", distance.Unix)
	}(time.Now())
	err = mw.next.Aggregate(ctx, distance)
	return
}

func (mw *loggingMiddleware) Calculate(ctx context.Context, id int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		if err != nil {
			mw.logger.Log("method", "Calculate", "error", err,
				"took", time.Since(start), "obuID", id)
		} else {
			mw.logger.Log("method", "Calculate", "error", err,
				"took", time.Since(start), "obuID", id, "totalDistance", invoice.TotalDistance,
				"totalAmount", invoice.TotalAmount)
		}
	}(time.Now())
	invoice, err = mw.next.Calculate(ctx, id)
	return
}

type instrumentationMiddleware struct {
	next   Service
	logger log.Logger
}

func newInstrumentationMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &instrumentationMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

func (mw *instrumentationMiddleware) Aggregate(ctx context.Context, distance types.Distance) error {

	return mw.next.Aggregate(ctx, distance)
}

func (mw *instrumentationMiddleware) Calculate(ctx context.Context, id int) (*types.Invoice, error) {
	return mw.next.Calculate(ctx, id)
}
