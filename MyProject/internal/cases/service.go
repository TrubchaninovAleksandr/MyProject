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
	GetAvgCoinRates(ctx context.Context, titles []string) ([]entities.Coin, error)
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

	coinsFromAPI, err := s.rateClient.GetCoinRates(ctx, titles)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных от API: %w", err)
	}

	if len(coinsFromAPI) == 0 {
		return []entities.Coin{}, nil
	}

	coinsFromDB, err := s.storage.GetCoins(ctx, titles)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных из хранилища: %w", err)
	}

	missingTitles := getMissingTitles(titles, coinsFromDB)
	if len(missingTitles) > 0 {
		coinsToStore := make([]entities.Coin, 0, len(coinsFromAPI))
		for _, coin := range coinsFromAPI {
			if containsTitle(missingTitles, coin.Title) {
				coinsToStore = append(coinsToStore, coin)
			}
		}

		if len(coinsToStore) > 0 {
			if err := s.storage.StoreCoins(ctx, coinsToStore); err != nil {
				return nil, fmt.Errorf("ошибка сохранения данных в хранилище: %w", err)
			}
		}
	}

	return orderCoinsByTitles(titles, coinsFromAPI), nil
}

func (s *Service) GetMax(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, fmt.Errorf("введите валюту")
	}

	coins, err := s.loadCoins(ctx, titles)
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

	coins, err := s.loadCoins(ctx, titles)
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
	if len(titles) == 0 {
		return nil, fmt.Errorf("введите валюту")
	}

	coins, err := s.rateClient.GetAvgCoinRates(ctx, titles)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения средних данных от API: %w", err)
	}

	if len(coins) == 0 {
		return []entities.Coin{}, nil
	}

	return orderCoinsByTitles(titles, coins), nil
}

// loadCoins получаем монеты из БД, догружаем недостающие из API и сохраняем их в БД.
func (s *Service) loadCoins(ctx context.Context, titles []string) ([]entities.Coin, error) {
	coinsFromDB, err := s.storage.GetCoins(ctx, titles)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных из хранилища: %w", err)
	}

	missingTitles := getMissingTitles(titles, coinsFromDB)
	if len(missingTitles) == 0 {
		return coinsFromDB, nil
	}

	coinsFromAPI, err := s.rateClient.GetCoinRates(ctx, missingTitles)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных от API: %w", err)
	}

	if len(coinsFromAPI) > 0 {
		if err := s.storage.StoreCoins(ctx, coinsFromAPI); err != nil {
			return nil, fmt.Errorf("ошибка сохранения данных в хранилище: %w", err)
		}
	}

	return append(coinsFromDB, coinsFromAPI...), nil
}

// getMissingTitles возвращает titles, которых нет в уже найденных монетах.
func getMissingTitles(titles []string, coins []entities.Coin) []string {
	missing := make([]string, 0)

	for _, title := range titles {
		found := false
		for _, coin := range coins {
			if coin.Title == title {
				found = true
				break
			}
		}

		if !found && !containsTitle(missing, title) {
			missing = append(missing, title)
		}
	}

	return missing
}

func containsTitle(titles []string, title string) bool {
	for _, current := range titles {
		if current == title {
			return true
		}
	}

	return false
}

// orderCoinsByTitles возвращает монеты в порядке входного titles.
func orderCoinsByTitles(titles []string, coins []entities.Coin) []entities.Coin {
	ordered := make([]entities.Coin, 0, len(titles))

	for _, title := range titles {
		for _, coin := range coins {
			if coin.Title == title {
				ordered = append(ordered, coin)
				break
			}
		}
	}

	return ordered
}
