package ports

import (
	"MyProject/internal/entities"
	"context"
)

type UserPort interface {
	GetMax(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetMin(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetAvg(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetPerc(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetLast(ctx context.Context, titles []string) ([]entities.Coin, error)
}
