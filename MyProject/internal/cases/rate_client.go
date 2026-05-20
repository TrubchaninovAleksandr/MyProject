package cases

import (
	"MyProject/internal/entities"
	"context"
)

// RateClient интерфейс для получения текущих курсов криптовалют из внешнего источника.

type RateClient interface {
	GetCoinRates(ctx context.Context, titles []string) ([]entities.Coin, error)
}
