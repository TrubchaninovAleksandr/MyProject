package adapters

import (
	"MyProject/internal/cases"
	"MyProject/internal/entities"
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

const (
	tableTitles = "titles" // столбец: title
	tableCoins  = "coins"  // столбцы: title, rate, creation_time
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

// NewPostgresStorage создаём пул подключений и возвращаем БД и закрываем
func NewPostgresStorage(ctx context.Context, dsn string) (*PostgresStorage, error) {
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}
	return &PostgresStorage{pool: pool}, nil
}

func (s *PostgresStorage) Close() {
	s.pool.Close()
}

// Query → открыть курсор

//rows.Next() → получить строку

//rows.Scan() → записать в переменную

//append → добавить в срез

//rows.Next() = false → выход из цикла

//rows.Err() → проверить ошибки внутри цикла

// rows.Close() → освободить соединение (через defer)
func (s *PostgresStorage) GetTitles(ctx context.Context) ([]string, error) {
	query, args, err := psql.
		Select("title").
		From(tableTitles).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ошибка построения запроса: %w", err)
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var titles []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		titles = append(titles, title)
	}

	return titles, rows.Err()
}

func (s *PostgresStorage) StoreCoins(ctx context.Context, coins []entities.Coin) error {
	if len(coins) == 0 {
		return nil
	}

	b := psql.Insert(tableCoins).Columns("title", "rate", "creation_time")
	now := time.Now()
	for _, c := range coins {
		b = b.Values(c.Title, c.Rate, now)
	}

	query, args, err := b.ToSql()
	if err != nil {
		return fmt.Errorf("ошибка построения запроса: %w", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %w", err)
	}

	return nil
}

func (s *PostgresStorage) GetCoins(ctx context.Context, titles []string, opts ...cases.CoinOption) ([]entities.Coin, error) {
	cfg := &cases.CoinConfig{}
	for _, o := range opts {
		o(cfg)
	}

	var query string
	var args []interface{}
	var err error

	switch cfg.Mode {
	case cases.Max:
		query, args, err = psql.
			Select("title", "MAX(rate) AS rate").
			From(tableCoins).
			Where(sq.Eq{"title": titles}).
			GroupBy("title").
			ToSql()
	case cases.Min:
		query, args, err = psql.
			Select("title", "MIN(rate) AS rate").
			From(tableCoins).
			Where(sq.Eq{"title": titles}).
			GroupBy("title").
			ToSql()
	case cases.Last:
		query, args, err = psql.
			Select("DISTINCT ON (title) title", "rate").
			From(tableCoins).
			Where(sq.Eq{"title": titles}).
			OrderBy("title", "creation_time DESC").
			ToSql()
	case cases.Perc:
		// Процент : (max - min) / min * 100 для каждой монеты
		query, args, err = psql.
			Select("title", "((MAX(rate) - MIN(rate)) / MIN(rate), 0)) * 100 AS rate").
			From(tableCoins).
			Where(sq.Eq{"title": titles}).
			GroupBy("title").
			ToSql()

	default:
		query, args, err = psql.
			Select("title", "rate").
			From(tableCoins).
			Where(sq.Eq{"title": titles}).
			OrderBy("creation_time DESC").
			ToSql()
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка построения запроса: %w", err)
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var result []entities.Coin
	for rows.Next() {
		var c entities.Coin
		if err := rows.Scan(&c.Title, &c.Rate); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		result = append(result, c)
	}

	return result, rows.Err()
}

func (s *PostgresStorage) GetAvgCoinsLastHour(ctx context.Context, titles []string) ([]entities.Coin, error) {
	query, args, err := psql.
		Select("title", "AVG(rate) AS rate").
		From(tableCoins).
		Where(sq.And{
			sq.Eq{"title": titles},
			sq.GtOrEq{"creation_time": time.Now().Add(-1 * time.Hour)},
		}).
		GroupBy("title").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("ошибка построения запроса: %w", err)
	}

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var coins []entities.Coin
	for rows.Next() {
		var c entities.Coin
		if err := rows.Scan(&c.Title, &c.Rate); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		coins = append(coins, c)
	}

	return coins, rows.Err()
}
