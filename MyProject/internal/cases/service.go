package cases

import (
	"MyProject/internal/entities"
	"context"
	"fmt"
)

type Service struct {
	rateClient RateClient
	storage    Storage
}

func NewService(rateClient RateClient, storage Storage) (*Service, error) {
	if rateClient == nil {
		return nil, fmt.Errorf("ошибка клиента для получения курсов ")
	}
	if storage == nil {
		return nil, fmt.Errorf("ошибка хранилища")
	}

	return &Service{
		rateClient: rateClient,
		storage:    storage,
	}, nil
}

type RateClient interface {
	GetCoinRates(ctx context.Context, titles []string) ([]entities.Coin, error)
}

type Storage interface {
	GetTitles(ctx context.Context) ([]string, error)
	StoreCoins(ctx context.Context, coins []entities.Coin) error
	GetCoins(ctx context.Context, titles []string) ([]entities.Coin, error)
}

func (s *Service) FetchRates(ctx context.Context) error {
	titles, err := s.storage.GetTitles(ctx)
	if err != nil {
		return err
	}

	if len(titles) == 0 {
		return nil
	}

	coins, err := s.rateClient.GetCoinRates(ctx, titles)
	if err != nil {
		return err
	}

	if len(coins) == 0 {
		return nil
	}

	if err := s.storage.StoreCoins(ctx, coins); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetActual(ctx context.Context, titles []string) ([]entities.Coin, error) {

	if len(titles) == 0 {
		return nil, fmt.Errorf("введите валюту")
	}

	coins, err := s.storage.GetCoins(ctx, titles)
	if err != nil {
		return nil, err
	}

	return coins, nil
}

// getLoadCoins получает монеты из БД, а если их нет — загружает из API и сохраняет
func (s *Service) getLoadCoins(ctx context.Context, titles []string) ([]entities.Coin, error) {

	coins, err := s.storage.GetCoins(ctx, titles)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных из хранилища: %w", err)
	}

	if len(coins) > 0 {
		return coins, nil
	}

	coins, err = s.rateClient.GetCoinRates(ctx, titles)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных от API: %w", err)
	}

	if len(coins) > 0 {
		if err := s.storage.StoreCoins(ctx, coins); err != nil {
			return nil, fmt.Errorf("ошибка сохранения данных в хранилище: %w", err)
		}
	}

	return coins, nil
}

func (s *Service) GetMax(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, fmt.Errorf("введите валюту")
	}

	// Получаем монеты (из БД или API)
	coins, err := s.getLoadCoins(ctx, titles)
	if err != nil {
		return nil, err
	}

	if len(coins) == 0 {
		return []entities.Coin{}, nil
	}

	maxCoin := coins[0]
	for _, coin := range coins[1:] {
		if coin.Rate > maxCoin.Rate {
			maxCoin = coin
		}
	}

	return []entities.Coin{maxCoin}, nil

}

func (s *Service) GetMin(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, fmt.Errorf("введите валюту")
	}

	coins, err := s.getLoadCoins(ctx, titles)
	if err != nil {
		return nil, err
	}

	if len(coins) == 0 {
		return []entities.Coin{}, nil
	}

	minCoin := coins[0]
	for _, coin := range coins[1:] {
		if coin.Rate < minCoin.Rate {
			minCoin = coin
		}
	}

	return []entities.Coin{minCoin}, nil
}

func (s *Service) GetAvg(ctx context.Context, titles []string) ([]entities.Coin, error) {
	// Получить курсы валют из хранилища
	// Вернуть курсы валют для конкретных валют из хранилища, и показать изменение их за последний час в процентах
	return []entities.Coin{}, nil
}
