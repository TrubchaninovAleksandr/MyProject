package ports

import (
	"MyProject/internal/entities"
	"context"
)

// UserPort интерфейс, описывающий контракт для получения курсов криптовалют.
type Service interface {
	GetMax(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetMin(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetAvg(ctx context.Context, titles []string) ([]entities.Coin, error)
	GetActual(ctx context.Context, titles []string) ([]entities.Coin, error)
}
