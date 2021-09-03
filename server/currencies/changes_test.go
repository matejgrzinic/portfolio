package currencies

import (
	"fmt"
	"testing"
	"time"

	"github.com/matejgrzinic/portfolio/db"
	"github.com/matejgrzinic/portfolio/external"
	"github.com/stretchr/testify/assert"
)

func Test_updateChangesMap(t *testing.T) {
	mockedDB := db.NewMockedDB()
	c := &Currencies{db: mockedDB}

	t.Run("unknown error ", func(t *testing.T) {
		mockedDB.QueryRowFunc = func(result interface{}) error {
			return fmt.Errorf("unittest")
		}

		err := c.updateChangesMap()
		assert.EqualError(t, err, "unittest")
	})

	t.Run("invalid data ", func(t *testing.T) {
		mockedDB.QueryRowFunc = func(result interface{}) error {
			return nil
		}

		err := c.updateChangesMap()
		assert.EqualError(t, err, "data field does not exist in dbData timeframe: hour")
	})

	t.Run("OK", func(t *testing.T) {
		mockFetcher := new(mockedfetcher)
		mockFetcher.fetchFunc = func() (external.CurrenciesDataMap, error) { return nil, nil }

		c.data = map[string]*currencyType{
			"fiat": newCurrecyType(mockFetcher, time.Minute),
		}
		c.data["fiat"].stopRefreshChan <- struct{}{}

		mockedDB.QueryRowFunc = func(result interface{}) error {
			result.(dbPriceData)["data"] = priceData{
				"fiat": map[string]external.CurrencyData{
					"USD": {
						Symbol: "USD",
						Price:  1.0,
					},
				},
			}
			return nil
		}

		err := c.updateChangesMap()
		assert.EqualError(t, err, "invalid symbol: USD")
	})

	t.Run("OK", func(t *testing.T) {
		mockFetcher := new(mockedfetcher)
		mockFetcher.fetchFunc = func() (external.CurrenciesDataMap, error) {
			return map[string]external.CurrencyData{
				"USD": {Symbol: "USD", Price: 1.0},
			}, nil
		}

		c.data = map[string]*currencyType{
			"fiat": newCurrecyType(mockFetcher, time.Minute),
		}
		c.data["fiat"].stopRefreshChan <- struct{}{}

		counter := 1.0
		mockedDB.QueryRowFunc = func(result interface{}) error {
			counter *= 2
			result.(dbPriceData)["data"] = priceData{
				"fiat": map[string]external.CurrencyData{
					"USD": {
						Symbol: "USD",
						Price:  counter,
					},
				},
			}
			return nil
		}

		err := c.updateChangesMap()
		assert.NoError(t, err)
		assert.Equal(t, -0.5, c.changes["fiat"]["USD"][Hour])
		assert.Equal(t, -0.75, c.changes["fiat"]["USD"][Day])
		assert.Equal(t, -0.875, c.changes["fiat"]["USD"][Week])
		assert.Equal(t, -0.9375, c.changes["fiat"]["USD"][Month])
	})

}
