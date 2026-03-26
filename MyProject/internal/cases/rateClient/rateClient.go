package rateClient

import (
	"MyProject/internal/entities"
	"context"
)

type RateClient interface {
	GetCoinRates(ctx context.Context, titles []string) ([]entities.Coin, error)
}
