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

// получение названия монет из БД
func (s *PostgresStorage) GetTitles(ctx context.Context) ([]string, error) {
}

// сохраняем монеты в БД
func (s *PostgresStorage) StoreCoins(ctx context.Context, coins []entities.Coin) error {

}

// получаем монеты по названиям из БД
func (s *PostgresStorage) GetCoins(ctx context.Context, titles []string) ([]entities.Coin, error) {

}
