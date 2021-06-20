package external

type CurrencyData struct {
	Symbol string  `json:"fiat"`
	Price  float64 `json:"cryptocurrency"`
}

type CurrenciesDataMap map[string]CurrencyData

type Fetcher interface {
	Fetch() (CurrenciesDataMap, error)
}
