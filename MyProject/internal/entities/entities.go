package entities

type Coin struct {
	Title string
	Rate  float64
}

func NewCoin(title string, rate float64) (*Coin, error) {

	return &Coin{
		Title: title,
		Rate:  rate,
	}, nil
}
