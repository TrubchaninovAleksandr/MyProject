package dto

type CoinsDTO []CoinDTO

type CoinDTO struct {
	Title string  `json:"title"`
	Rate  float64 `json:"rate"`
}
