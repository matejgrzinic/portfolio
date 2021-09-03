package currencies

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/matejgrzinic/portfolio/db"
	"github.com/matejgrzinic/portfolio/external"
	"github.com/matejgrzinic/portfolio/refresher"
)

var (
	CryptocurrencyType = "cryptocurrency"
	FiatType           = "fiat"
)

type Currencies struct {
	data    map[string]*currencyType
	changes changesMap
	db      db.API

	dbSaveInterval time.Duration

	resultSaveChan chan error
	stopSaveChan   chan struct{}
}

type CTX interface {
	DB() db.API
}

func NewCurrencies(ctx CTX) *Currencies {
	c := new(Currencies)

	c.data = map[string]*currencyType{
		FiatType:           newCurrecyType(external.NewFiatFetcher(), time.Hour),
		CryptocurrencyType: newCurrecyType(external.NewCryptocurrencyFetcher(), time.Minute),
	}
	c.db = ctx.DB()

	dbIntervalStr := os.Getenv("DB_SAVE_INTERVAL_SECONDS")
	dbInterval, err := time.ParseDuration(dbIntervalStr)
	if err != nil {
		log.Panicf("invalid DB_SAVE_INTERVAL_SECONDS env variable: %v", dbIntervalStr)
	}
	c.dbSaveInterval = dbInterval
	c.stopSaveChan = make(chan struct{}, 1)
	c.resultSaveChan = make(chan error, 1)

	refresher.StartRefresher(c.resultSaveChan, c.stopSaveChan, c.dbSaveInterval, c.saveToDbAndUpdateChanges)
	<-c.resultSaveChan

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
	Changes CurrencyChanges `json:"changes"`
}

func (c *Currencies) GetCurrencyWithChanges(currencyType, symbol string) (*CurrencyDataWithChanges, error) {
	currencyData, err := c.GetCurrency(currencyType, symbol)
	if err != nil {
		return nil, err
	}

	changes, err := c.getChangesForCurrency(currencyType, symbol)
	if err != nil {
		return nil, err
	}

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
