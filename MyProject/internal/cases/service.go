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

func (s *Service) FetchRates(ctx context.Context) error {
	titles, err := s.storage.GetTitles(ctx)
	if err != nil {
		return fmt.Errorf("ошибка получения списка валют из хранилища: %w", err)
	}

	if len(titles) == 0 {
		return fmt.Errorf("список валют в хранилище пуст")
	}

	coins, err := s.rateClient.GetCoinRates(ctx, titles)
	if err != nil {
		return fmt.Errorf("ошибка получения курсов от API: %w", err)
	}

	if len(coins) == 0 {
		return fmt.Errorf("API вернул пустой список курсов для валют: %v", titles)
	}

	if err := s.storage.StoreCoins(ctx, coins); err != nil {
		return fmt.Errorf("ошибка сохранения курсов в хранилище: %w", err)
	}

	return nil
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

	if err := s.missingTitles(ctx, titles); err != nil {
		return nil, fmt.Errorf("ошибка обработки недостающих валют: %w", err)
	}

	coins, err := s.storage.GetCoins(ctx, titles, CoinMax())
	if err != nil {
		return nil, fmt.Errorf("ошибка получения максимальных курсов из хранилища: %w", err)
	}

	return coins, nil
}

func (s *Service) GetMin(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, fmt.Errorf("введите валюту")
	}

	if err := s.missingTitles(ctx, titles); err != nil {
		return nil, fmt.Errorf("ошибка обработки недостающих валют: %w", err)
	}

	coins, err := s.storage.GetCoins(ctx, titles, CoinMin())
	if err != nil {
		return nil, fmt.Errorf("ошибка получения минимальных курсов из хранилища: %w", err)
	}

	return coins, nil
}

func (s *Service) GetAvg(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, fmt.Errorf("введите валюту")
	}

	if err := s.missingTitles(ctx, titles); err != nil {
		return nil, fmt.Errorf("ошибка обработки недостающих валют: %w", err)
	}

	coins, err := s.storage.GetAvg(ctx, titles)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения средних курсов из хранилища: %w", err)
	}

	return coins, nil
}

func (s *Service) GetActual(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, fmt.Errorf("введите валюту")
	}

	if err := s.missingTitles(ctx, titles); err != nil {
		return nil, fmt.Errorf("ошибка обработки недостающих валют: %w", err)
	}

	coins, err := s.storage.GetCoins(ctx, titles, CoinActual())
	if err != nil {
		return nil, fmt.Errorf("ошибка получения последних курсов из хранилища: %w", err)
	}

	return coins, nil
}
