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

// CoinOption тип функции для применения опций конфигурации.
type CoinOption func(cf *CoinConfig)

// CoinMax возвращает опцию для получения максимального курса.
func CoinMax() CoinOption {
	return func(cf *CoinConfig) {
		cf.Mode = Max
	}
}

// CoinMin возвращает опцию для получения минимального курса.
func CoinMin() CoinOption {
	return func(cf *CoinConfig) {
		cf.Mode = Min
	}
}

// CoinPerc возвращает опцию для получения процента изменения курса.
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

// String возвращает строковое представление режима CoinMode.
func (c CoinMode) String() string {
	return [...]string{"", "Max", "Min", "Perc", "Last"}[c]
}
