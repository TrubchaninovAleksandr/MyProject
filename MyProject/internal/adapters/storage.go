package adapters

import (
	"MyProject/internal/entities"
	"context"
	"database/sql"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) GetTitles(ctx context.Context) ([]string, error) {
}

func (s *PostgresStorage) StoreCoins(ctx context.Context, coins []entities.Coin) error {

	return nil
}

func (s *PostgresStorage) GetCoins(ctx context.Context, titles []string) ([]entities.Coin, error) {
}
