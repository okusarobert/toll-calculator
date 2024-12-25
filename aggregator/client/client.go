package client

import (
	"context"

	"github.com/okusarobert/toll-calculator/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
}
