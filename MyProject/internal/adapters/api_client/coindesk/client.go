package coindesk

import (
	"MyProject/internal/entities"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL  = "https://min-api.cryptocompare.com"
	timeout  = 10 * time.Second
	currency = "USD"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	convert    string
}

func NewClient(baseURL string, timeoutSeconds int, currency string) (*Client, error) {
	timeout := time.Duration(timeoutSeconds) * time.Second

	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    baseURL,
		convert:    currency,
	}, nil
}

// GET /data/price?fsym=BTC&tsyms=USD

func (c *Client) GetCoinRates(ctx context.Context, titles []string) ([]entities.Coin, error) {
	if len(titles) == 0 {
		return nil, fmt.Errorf("список монет пуст")
	}

	coins := make([]entities.Coin, 0, len(titles))

	// GET /data/pricemulti?fsyms=BTC,ETH&tsyms=USD — один запрос для всех монет
	u, err := url.Parse(c.baseURL + "/data/pricemulti")
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга URL: %w", err)
	}
	q := u.Query()
	q.Set("fsyms", strings.Join(titles, ","))
	q.Set("tsyms", c.convert)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус %d: %s", resp.StatusCode, string(body))
	}

	// Ответ: {"BTC":{"USD":12345.67},"ETH":{"USD":456.78}}
	var result map[string]map[string]float64
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	for _, title := range titles {
		rates, ok := result[title]
		if !ok {
			return nil, fmt.Errorf("в ответе отсутствует монета %s", title)
		}
		rate, ok := rates[c.convert]
		if !ok {
			return nil, fmt.Errorf("в ответе отсутствует курс %s для %s", c.convert, title)
		}
		coin, err := entities.NewCoin(title, rate)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания монеты %s: %w", title, err)
		}
		coins = append(coins, *coin)
	}

	return coins, nil
}
