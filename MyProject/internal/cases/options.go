package cases

// CoinMode — тип режима сравнения.
type CoinMode int

// Режимы сравнения.
const (
	_ CoinMode = iota
	Max
	Min
	Avg
)

// CoinConfig — структура конфигурации выборки монет.
type CoinConfig struct {
	Mode CoinMode
}

// CoinOption тип функции для применения опций конфигурации.
type CoinOption func(cfg *CoinConfig)

// CoinMax возвращает опцию для получения максимального курса.
func CoinMax() CoinOption {
	return func(cfg *CoinConfig) {
		cfg.Mode = Max
	}
}

// CoinMin возвращает опцию для получения минимального курса.
func CoinMin() CoinOption {
	return func(cfg *CoinConfig) {
		cfg.Mode = Min
	}
}

// CoinAvg возвращает опцию для получения среднего курса.
func CoinAvg() CoinOption {
	return func(cfg *CoinConfig) {
		cfg.Mode = Avg
	}
}
