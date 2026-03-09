package entities

import "fmt"

type Coin struct {
	Title string
	Rate  float64
}

func NewCoin(title string, rate float64) (*Coin, error) {
	if title == "" {
		return nil, fmt.Errorf("название не может быть пустым")
	}
	if rate <= 0 {
		return nil, fmt.Errorf("курс должен быть положительным, получено: %f", rate)
	}

	return &Coin{
		Title: title,
		Rate:  rate,
	}, nil
}

//
