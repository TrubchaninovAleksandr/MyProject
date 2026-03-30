package service

type CoinOptions func(a, b float64) bool

func Max() CoinOptions { return func(a, b float64) bool { return a > b } }
func Min() CoinOptions { return func(a, b float64) bool { return a < b } }
