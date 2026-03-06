package cases

import (
	"awesomeProject/internal/entities"
	"context"

)

type Service struct {
	rateClient RateClient
	storage  Storage
}

func NewService(rateClient RateClient, storage Storage) (*Service, error) {
	return &Service{
		rateClient: rateClient,
		storage:    storage,
	}
}

type RateСlient interface {
	GetCoinRates(ctx context.Context, titles []string) ([]entities.Coin, error)
}


func (s *Service) FetchRates (ctx context.Context) error {
	titles, err := s.Storage.GetTitles(ctx)
	if err != nil {
		return err
	}

	if len(titles) == 0 {
		return nil
	}

	coins, err := s.rateClient.FetchRates(ctx, titles)
	if err != nil {
		return err
	}

	if len(coins) == 0 {
		return nil
	}

	if err := s.Storage.StoreCoins(ctx, coins); err != nil {
		return err
	}
	return nil
}

func checkTitles (ctx context.Context, titles []string) (_, error) {

}

func (s *Service) GetActual(ctx context.Context, titles []string) ([]entities.Coin, error) {
	// Получить курсы валют для конкретных валют из БД
	var result []entities.Coin

	return result, nil
}

func (s *Service) GetMax (ctx context.Context, titles []string) ([]entities.Coin, error) {
	// Получить курсы валют  из хранилища
	// Вернуть курс валюты с максимальной стоимостью для конкретной валюты из хранилища
	//Если нет в хранилище, то получить новые курсы валют от API и сохранить их в хранилище
	return entities.Coin, nil
}

func (s *Service) GetMin(ctx context.Context, titles []string) ([]entities.Coin, error) {
	// Получить курсы валют  из хранилища
	// Вернуть курс валюты с минимальной стоимостью для конкретных валют из хранилища
	//Если нет в хранилище, то получить новые курсы валют от API и сохранить их в хранилище
	return entities.Coin, nil
}

func (s *Service) GetAvg (ctx context.Context, titles []string) ([]entities.Coin, error) {
	// Получить курсы валют  из хранилища
	// Вернуть курсы валют для конкретных валют из хранилища, и показать изменение их за последний час в процентах