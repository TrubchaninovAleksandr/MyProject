package adapters

import (
	"MyProject/internal/entities"
	"context"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c Client) GetCoinRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	return nil, nil
}
