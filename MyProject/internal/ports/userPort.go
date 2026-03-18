package ports

import (
	"MyProject/internal/entities"
	"context"
)

type UserPort interface {
	GetActual(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetMax(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetMin(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetAvg(ctx context.Context, titles []string) ([]entities.Coin, error)
}
