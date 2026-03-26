package cases

import (
	"MyProject/internal/cases/rateClient"
	"MyProject/internal/cases/storage"
	"MyProject/internal/entities"
	"context"
	"fmt"
)

type Service struct {
	rateClient rateClient.RateClient
	storage    storage.Storage
}

func NewService(rateClient rateClient.RateClient, storage storage.Storage) (*Service, error) {
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

func (s *Service) FetchRates(ctx context.Context) error {
	titles, err := s.storage.GetTitles(ctx)
	if err != nil {
		return fmt.Errorf("ошибка получения списка валют из хранилища: %w", err)
	}

	if len(titles) == 0 {
		return nil
	}

	coins, err := s.rateClient.GetCoinRates(ctx, titles)
	if err != nil {
		return fmt.Errorf("ошибка получения курсов от API: %w", err)
	}

	if len(coins) == 0 {
		return nil
	}

	if err := s.storage.StoreCoins(ctx, coins); err != nil {
		return fmt.Errorf("ошибка сохранения курсов в хранилище: %w", err)
	}

	return nil
}

func (s *Service) GetActual(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, fmt.Errorf("введите валюту")
	}

	if err := s.missingTitles(ctx, titles); err != nil {
		return nil, fmt.Errorf("ошибка обработки недостающих валют: %w", err)
	}

	coins, err := s.storage.GetCoins(ctx, titles)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных из хранилища: %w", err)
	}

	return coins, nil
}

func (s *Service) missingTitles(ctx context.Context, titles []string) error {
	dbTitles, err := s.storage.GetTitles(ctx)
	if err != nil {
		return fmt.Errorf("ошибка получения списка валют из хранилища: %w", err)
	}

	missingTitles := make([]string, 0)
	for _, t := range titles {
		found := false
		for _, dt := range dbTitles {
			if dt == t {
				found = true
				break
			}
		}
		if !found {
			missingTitles = append(missingTitles, t)
		}
	}

	if len(missingTitles) == 0 {
		return nil
	}

	coinsFromAPI, err := s.rateClient.GetCoinRates(ctx, missingTitles)
	if err != nil {
		return fmt.Errorf("ошибка получения данных от API: %w", err)
	}

	if len(coinsFromAPI) > 0 {
		if err := s.storage.StoreCoins(ctx, coinsFromAPI); err != nil {
			return fmt.Errorf("ошибка сохранения данных в хранилище: %w", err)
		}
	}

	return nil
}

func (s *Service) GetMax(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, fmt.Errorf("введите валюту")
	}
	// Убедимся, что все переданные titles присутствуют в хранилище
	if err := s.missingTitles(ctx, titles); err != nil {
		return nil, err
	}

	// Теперь читаем все требуемые монеты из БД
	coins, err := s.storage.GetCoins(ctx, titles)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных из хранилища: %w", err)
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
	// Убедимся, что все переданные titles присутствуют в хранилище
	if err := s.missingTitles(ctx, titles); err != nil {
		return nil, err
	}

	// читаем монеты из БД
	coins, err := s.storage.GetCoins(ctx, titles)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных из хранилища: %w", err)
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
	if len(titles) == 0 {
		return nil, fmt.Errorf("введите валюту")
	}

	coins, err := s.storage.GetAvgCoinsLastHour(ctx, titles)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения средних данных из хранилища: %w", err)
	}

	return coins, nil
}
