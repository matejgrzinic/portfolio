package currencies

import "github.com/matejgrzinic/portfolio/external"

type API interface {
	GetCurrency(currencyType, symbol string) (*external.CurrencyData, error)
}
