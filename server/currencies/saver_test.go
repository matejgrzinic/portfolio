package currencies

import (
	"fmt"
	"testing"
	"time"

	"github.com/matejgrzinic/portfolio/db"
	"github.com/matejgrzinic/portfolio/external"
	"github.com/stretchr/testify/assert"
)

func Test_saveToDb(t *testing.T) {
	mockFetcher := new(mockedfetcher)
	mockedDB := db.NewMockedDB()
	mockFetcher2 := new(mockedfetcher)

	t.Run("convertToEur returns error", func(t *testing.T) {
		mockFetcher2.fetchFunc = func() (external.CurrenciesDataMap, error) {
			return map[string]external.CurrencyData{
				"BTC": {Symbol: "BTC", Price: 15000},
			}, nil
		}

		c := &Currencies{db: mockedDB}
		c.data = map[string]*currencyType{
			"cryptocurrency": newCurrecyType(mockFetcher2, time.Minute),
		}
		c.data["cryptocurrency"].stopRefreshChan <- struct{}{}

		mockedDB.InsertOneFunc = func(data interface{}) error {
			mydata, ok := data.(*struct {
				Data map[string]external.CurrenciesDataMap `json:"data"`
				Time int64                                 `json:"time"`
			})
			assert.True(t, ok)
			assert.Equal(t, 1.5, mydata.Data["fiat"]["USD"].Price)
			assert.Equal(t, 10000.0, mydata.Data["cryptocurrency"]["BTC"].Price)
			return nil
		}

		err := c.saveToDB()
		assert.EqualError(t, err, "save to db: fiat currency type does not exist")
	})

	t.Run("insertOne returns error", func(t *testing.T) {
		mockFetcher.fetchFunc = func() (external.CurrenciesDataMap, error) {
			return map[string]external.CurrencyData{
				"USD": {Symbol: "USD", Price: 1.5},
			}, nil
		}

		mockFetcher2.fetchFunc = func() (external.CurrenciesDataMap, error) {
			return map[string]external.CurrencyData{
				"BTC": {Symbol: "BTC", Price: 15000},
			}, nil
		}

		c := &Currencies{db: mockedDB}
		c.data = map[string]*currencyType{
			"fiat":           newCurrecyType(mockFetcher, time.Minute),
			"cryptocurrency": newCurrecyType(mockFetcher2, time.Minute),
		}
		c.data["fiat"].stopRefreshChan <- struct{}{}
		c.data["cryptocurrency"].stopRefreshChan <- struct{}{}

		mockedDB.InsertOneFunc = func(data interface{}) error {
			mydata, ok := data.(*struct {
				Data map[string]external.CurrenciesDataMap `json:"data"`
				Time int64                                 `json:"time"`
			})
			assert.True(t, ok)
			assert.Equal(t, 1.5, mydata.Data["fiat"]["USD"].Price)
			assert.Equal(t, 10000.0, mydata.Data["cryptocurrency"]["BTC"].Price)
			return fmt.Errorf("unittest")
		}

		err := c.saveToDB()
		assert.EqualError(t, err, "unittest")
	})

	t.Run("OK", func(t *testing.T) {
		mockFetcher.fetchFunc = func() (external.CurrenciesDataMap, error) {
			return map[string]external.CurrencyData{
				"USD": {Symbol: "USD", Price: 1.5},
			}, nil
		}

		mockFetcher2.fetchFunc = func() (external.CurrenciesDataMap, error) {
			return map[string]external.CurrencyData{
				"BTC": {Symbol: "BTC", Price: 15000},
			}, nil
		}

		c := &Currencies{db: mockedDB}
		c.data = map[string]*currencyType{
			"fiat":           newCurrecyType(mockFetcher, time.Minute),
			"cryptocurrency": newCurrecyType(mockFetcher2, time.Minute),
		}
		c.data["fiat"].stopRefreshChan <- struct{}{}
		c.data["cryptocurrency"].stopRefreshChan <- struct{}{}

		mockedDB.InsertOneFunc = func(data interface{}) error {
			mydata, ok := data.(*struct {
				Data map[string]external.CurrenciesDataMap `json:"data"`
				Time int64                                 `json:"time"`
			})
			assert.True(t, ok)
			assert.Equal(t, 1.5, mydata.Data["fiat"]["USD"].Price)
			assert.Equal(t, 10000.0, mydata.Data["cryptocurrency"]["BTC"].Price)
			return nil
		}

		err := c.saveToDB()
		assert.NoError(t, err)
	})
}
