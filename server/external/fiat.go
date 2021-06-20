package external

import (
	"encoding/json"
	"fmt"
)

type FiatFetcher struct {
	Name   string
	getAPI getAPI
}

func NewFiatFetcher() *FiatFetcher {
	f := new(FiatFetcher)
	f.Name = "fiat"
	f.getAPI = new(getImpl)
	return f
}

type fiatApiData struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}

func (f *FiatFetcher) Fetch() (CurrenciesDataMap, error) {
	data, err := f.getAPI.getFiat()
	if err != nil {
		return nil, err
	}

	var fd fiatApiData
	err = json.Unmarshal(data, &fd)
	if err != nil {
		return nil, fmt.Errorf("invalid data format for fiat api response: %v", err)
	}

	refreshedData := make(CurrenciesDataMap)
	for sym, price := range fd.Rates {
		refreshedData[sym] = CurrencyData{Symbol: sym, Price: price}
	}

	return refreshedData, nil
}
