package options

// CoinMode — тип режима сравнения.
type CoinMode int

// Режимы сравнения.
const (
	_ CoinMode = iota
	Max
	Min
)

// CoinConfig — структура конфигурации выборки монет.
type CoinConfig struct {
	Mode CoinMode
}

// CoinOption — функциональная опция для настройки CoinConfig.
type CoinOption func(*CoinConfig)

// Max — опция выбора монеты с максимальным курсом.
func CoinMax() CoinOption {
	return func(cfg *CoinConfig) {
		cfg.Mode = Max
	}
}

// Min — опция выбора монеты с минимальным курсом.
func CoinMin() CoinOption {
	return func(cfg *CoinConfig) {
		cfg.Mode = Min
	}
}

// NewCoinConfig — конструктор
func NewCoinConfig(opts ...CoinOption) *CoinConfig {
	cfg := &CoinConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// Compare сравнивает два значения согласно выбранной опции
func (cfg *CoinConfig) Compare(a, b float64) bool {
	switch cfg.Mode {
	case Max:
		return a > b
	case Min:
		return a < b
	default:
		return false
	}
}
