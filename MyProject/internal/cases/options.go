package cases

// CoinMode — тип режима сравнения.
type CoinMode int

// Режимы сравнения.
const (
	_ CoinMode = iota
	Max
	Min
	Perc
	Last
)

// CoinConfig — структура конфигурации выборки монет.
type CoinConfig struct {
	Mode CoinMode
}

type CoinOption func(cf *CoinConfig)

func CoinMax() CoinOption {
	return func(cf *CoinConfig) {
		cf.Mode = Max
	}
}

func CoinMin() CoinOption {
	return func(cf *CoinConfig) {
		cf.Mode = Min
	}
}
func CoinPerc() CoinOption {
	return func(cf *CoinConfig) {
		cf.Mode = Perc
	}
}

// CoinLast — опция выбора последнего самого свежего курса монеты.
func CoinLast() CoinOption {
	return func(cf *CoinConfig) {
		cf.Mode = Last
	}
}

func (c CoinMode) String() string {
	return [...]string{"", "Max", "Min", "Perc", "Last"}[c]
}
