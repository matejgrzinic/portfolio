package external

type CurrencyData struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

type CurrenciesDataMap map[string]CurrencyData

type Fetcher interface {
	Fetch() (CurrenciesDataMap, error)
}
