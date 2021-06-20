package external

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type CryptocurrencyFetcher struct {
	Name        string
	externalApi getAPI
}

func NewCryptocurrencyFetcher() *CryptocurrencyFetcher {
	c := new(CryptocurrencyFetcher)
	c.Name = "cryptocurrency"
	c.externalApi = new(getImpl)
	return c
}

type cryptoApiData struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func (c *CryptocurrencyFetcher) Fetch() (CurrenciesDataMap, error) {
	data, err := c.externalApi.getCryptocurrency()
	if err != nil {
		return nil, err
	}

	var cd []cryptoApiData
	err = json.Unmarshal(data, &cd)
	if err != nil {
		return nil, fmt.Errorf("invalid data format for cryptocurrency api response: %v", err)
	}

	refreshedData := make(CurrenciesDataMap)
	for _, row := range cd {
		if strings.HasSuffix(row.Symbol, "USDT") {
			symbol := row.Symbol[:len(row.Symbol)-4]
			priceFloat64, err := strconv.ParseFloat(row.Price, 64)
			if err != nil {
				return nil, fmt.Errorf("can not parse float from cryptocurrency price: (symbol:%v) (value:%v): %v", symbol, row.Price, err)
			}
			refreshedData[symbol] = CurrencyData{Symbol: symbol, Price: priceFloat64}
		}
	}

	return refreshedData, nil
}
