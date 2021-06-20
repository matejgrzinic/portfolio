package currencies

import (
	"fmt"
	"time"

	"github.com/matejgrzinic/portfolio/db"
	"github.com/matejgrzinic/portfolio/external"
)

var (
	CryptocurrencyType = "cryptocurrency"
	FiatType           = "fiat"
)

type changesMap map[string]map[string]map[string]float64

type Currencies struct {
	data    map[string]*currencyType
	changes changesMap
	dbAPI   *db.DB
	// prices handler?? // save prices + calculate changes
	// readyChan chan struct{} // currency type signal when data is rdy first time
}

func NewCurrencies() *Currencies {
	c := new(Currencies)

	c.data = map[string]*currencyType{
		FiatType:           newCurrecyType(external.NewCryptocurrencyFetcher(), time.Hour),
		CryptocurrencyType: newCurrecyType(external.NewFiatFetcher(), time.Minute),
	}

	return c
}

func (c *Currencies) GetCurrency(currencyType, symbol string) (*external.CurrencyData, error) {
	var data *external.CurrencyData
	var err error
	if dataType, ok := c.data[currencyType]; ok {
		data, err = dataType.Get(symbol)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("invalid currency type: %v", currencyType)
	}

	dataOut := &external.CurrencyData{Symbol: data.Symbol, Price: data.Price}
	err = c.convertToEur(dataOut, currencyType)
	if err != nil {
		return nil, err
	}
	return dataOut, nil
}

type CurrencyDataWithChanges struct {
	external.CurrencyData
	Changes map[string]float64
}

// func (c *Currencies) getChangesForCurrency(currencyType, symbol string) {

// }

func (c *Currencies) GetCurrencyWithChanges(currencyType, symbol string) (*CurrencyDataWithChanges, error) {
	currencyData, err := c.GetCurrency(currencyType, symbol)
	if err != nil {
		return nil, err
	}
	changes := map[string]float64{"hour": 1.04, "day": 1.10}

	// getChanges(type, symbol) map, err
	// -> cached nekje
	// --> from prices where time > now-year
	// --> refresh every prices insert
	// ---> insert every 10min?

	cwc := CurrencyDataWithChanges{CurrencyData: *currencyData, Changes: changes}
	return &cwc, nil
}

func (c *Currencies) convertToEur(curr *external.CurrencyData, currencyType string) error {
	if currencyType != FiatType {
		if _, ok := c.data[FiatType]; !ok {
			return fmt.Errorf("fiat currency type does not exist")
		}
		eurUsd, err := c.data[FiatType].Get("USD")
		if err != nil {
			return err
		}
		curr.Price = curr.Price / eurUsd.Price
	}
	return nil
}
