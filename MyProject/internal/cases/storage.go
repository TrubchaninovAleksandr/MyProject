package cases

import (
	"MyProject/internal/entities"
	"context"
)

type Storage interface {
	GetTitles(ctx context.Context) ([]string, error)
	StoreTitles(ctx context.Context, titles []string) error
	StoreCoins(ctx context.Context, coins []entities.Coin) error
	GetCoins(ctx context.Context, titles []string, opt ...CoinOption) ([]entities.Coin, error)
}
