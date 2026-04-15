package adapters

import (
	"MyProject/internal/entities"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	baseURL  = "https://min-api.cryptocompare.com"
	timeout  = 10 * time.Second
	currency = "USD"
	apiKey   = "2466ecb1d8ac4f9ecaf1d81c8370c32ab7ee9d920ac8b395617b91fba34cfaf2"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	convert    string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    baseURL,
		convert:    currency,
		apiKey:     apiKey,
	}
}

// GET /data/price?fsym=BTC&tsyms=USD

func (c *Client) GetCoinRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, fmt.Errorf("список монет пуст")
	}

	coins := make([]entities.Coin, 0, len(titles))

	for _, title := range titles {
		url := fmt.Sprintf("%s/data/price?fsym=%s&tsyms=%s", c.baseURL, title, c.convert)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания запроса для %s: %w", title, err)
		}

		req.Header.Set("Authorization", "Apikey "+c.apiKey)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("ошибка выполнения запроса для %s: %w", title, err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("неожиданный статус %d для %s: %s", resp.StatusCode, title, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения ответа для %s: %w", title, err)
		}

		var result map[string]float64
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("ошибка парсинга JSON для %s: %w", title, err)
		}

		rate, ok := result[c.convert]
		if !ok {
			return nil, fmt.Errorf("в ответе отсутствует ключ %s для %s: %s", c.convert, title, string(body))
		}

		coin, err := entities.NewCoin(title, rate)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания монеты %s: %w", title, err)
		}
		coins = append(coins, *coin)
	}

	return coins, nil
}
