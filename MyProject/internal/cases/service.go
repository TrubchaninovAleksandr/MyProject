package cases

import (
	service "MyProject/internal/cases/options"
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

func (s *Service) getCoinByOption(ctx context.Context, titles []string, opt service.CoinOptions) ([]entities.Coin, error) {
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
	if len(coins) == 0 {
		return []entities.Coin{}, nil
	}
	best := coins[0]
	for _, c := range coins[1:] {
		if opt(c.Rate, best.Rate) {
			best = c
		}
	}
	return []entities.Coin{best}, nil
}

func (s *Service) GetMax(ctx context.Context, titles []string) ([]entities.Coin, error) {
	return s.getCoinByOption(ctx, titles, service.Max())
}

func (s *Service) GetMin(ctx context.Context, titles []string) ([]entities.Coin, error) {
	return s.getCoinByOption(ctx, titles, service.Min())
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
