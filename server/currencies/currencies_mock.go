package currencies

import "github.com/matejgrzinic/portfolio/external"

type MockedCurrencies struct {
	GetCurrencyFunc            func(currencyType, symbol string) (*external.CurrencyData, error)
	GetCurrencyWithChangesFunc func(currencyType, symbol string) (*CurrencyDataWithChanges, error)
}

func NewMockedCurrencies() *MockedCurrencies {
	return new(MockedCurrencies)
}

func (mc *MockedCurrencies) GetCurrency(currencyType, symbol string) (*external.CurrencyData, error) {
	return mc.GetCurrencyFunc(currencyType, symbol)
}

func (mc *MockedCurrencies) GetCurrencyWithChanges(currencyType, symbol string) (*CurrencyDataWithChanges, error) {
	return mc.GetCurrencyWithChangesFunc(currencyType, symbol)
}
