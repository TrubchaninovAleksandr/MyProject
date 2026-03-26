package entities

import "fmt"

type Coin struct {
	Title string
	Rate  float64
}

func NewCoin(title string, rate float64) (*Coin, error) {
	if title == "" {
		return nil, fmt.Errorf("title не может быть пустым")
	}
	if rate < 0 {
		return nil, fmt.Errorf("rate не может быть отрицательным: %v", rate)
	}

	return &Coin{
		Title: title,
		Rate:  rate,
	}, nil
}
