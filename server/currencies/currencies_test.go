package currencies

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/matejgrzinic/portfolio/external"

	"github.com/stretchr/testify/assert"
)

type mockedfetcher struct {
	fetchFunc func() (external.CurrenciesDataMap, error)
}

func (f *mockedfetcher) Fetch() (external.CurrenciesDataMap, error) {
	return f.fetchFunc()
}

func TestGetCurrency(t *testing.T) {
	mockFetcher := new(mockedfetcher)
	mockFetcher.fetchFunc = func() (external.CurrenciesDataMap, error) {
		return map[string]external.CurrencyData{
			"USD": {Symbol: "USD", Price: 2.0},
		}, nil
	}

	c := new(Currencies)
	c.data = map[string]*currencyType{
		"unittest": newCurrecyType(mockFetcher, time.Minute),
	}
	c.data["unittest"].stopRefreshChan <- struct{}{}

	t.Run("invalid symbol", func(t *testing.T) {
		curr, err := c.GetCurrency("unittest", "invalid")
		assert.Nil(t, curr)
		assert.EqualError(t, err, "invalid symbol: invalid")
	})

	t.Run("invalid currency type", func(t *testing.T) {
		curr, err := c.GetCurrency("unittest_invalid", "invalid")
		assert.Nil(t, curr)
		assert.EqualError(t, err, "invalid currency type: unittest_invalid")
	})

	t.Run("invalid currency type", func(t *testing.T) {
		curr, err := c.GetCurrency("unittest_invalid", "invalid")
		assert.Nil(t, curr)
		assert.EqualError(t, err, "invalid currency type: unittest_invalid")
	})

	t.Run("fiat does not exist for conversion", func(t *testing.T) {
		curr, err := c.GetCurrency("unittest", "USD")
		assert.Nil(t, curr)
		assert.EqualError(t, err, "fiat currency type does not exist")
	})

	c.data["fiat"] = newCurrecyType(mockFetcher, time.Minute)
	c.data["fiat"].stopRefreshChan <- struct{}{}

	t.Run("fiat get USD returns error", func(t *testing.T) {
		c.data["fiat"].data["USDCOPY"] = c.data["fiat"].data["USD"]
		delete(c.data["fiat"].data, "USD")
		defer func() {
			c.data["fiat"].data["USD"] = c.data["fiat"].data["USDCOPY"]
			delete(c.data["fiat"].data, "USDCOPY")
		}()

		curr, err := c.GetCurrency("unittest", "USD")
		assert.Nil(t, curr)
		assert.EqualError(t, err, "invalid symbol: USD")
	})

	t.Run("OK fiat", func(t *testing.T) {
		curr, err := c.GetCurrency("fiat", "USD")
		assert.NoError(t, err)
		assert.Equal(t, curr.Symbol, "USD")
		assert.Equal(t, curr.Price, 2.0)
		assert.Equal(t, c.data["fiat"].data["USD"].Price, 2.0)
	})

	t.Run("OK not fiat", func(t *testing.T) {
		curr, err := c.GetCurrency("unittest", "USD")
		assert.NoError(t, err)
		assert.Equal(t, curr.Symbol, "USD")
		assert.Equal(t, curr.Price, 1.0)
		assert.Equal(t, c.data["unittest"].data["USD"].Price, 2.0)
	})
}

func TestGetCurrencyWithChanges(t *testing.T) {
	mockFetcher := new(mockedfetcher)
	mockFetcher.fetchFunc = func() (external.CurrenciesDataMap, error) {
		return map[string]external.CurrencyData{
			"USD": {Symbol: "USD", Price: 2.0},
		}, nil
	}

	c := new(Currencies)
	c.data = map[string]*currencyType{
		"fiat": newCurrecyType(mockFetcher, time.Minute),
	}
	c.data["fiat"].stopRefreshChan <- struct{}{}

	log.Println(c.GetCurrencyWithChanges("fiat", "USD"))
}

func TestRefresh(t *testing.T) {
	mockFetcher := new(mockedfetcher)

	mockFetcher.fetchFunc = func() (external.CurrenciesDataMap, error) {
		return nil, fmt.Errorf("some_error")
	}
	ct := newCurrecyType(mockFetcher, time.Minute)
	ct.stopRefreshChan <- struct{}{}
	assert.Equal(t, ct.data, external.CurrenciesDataMap{})

	mockFetcher.fetchFunc = func() (external.CurrenciesDataMap, error) {
		return map[string]external.CurrencyData{
			"USD": {Symbol: "USD", Price: 0.0},
		}, nil
	}

	ct.refreshInterval = time.Millisecond * 100
	go ct.startRefresher()
	mockFetcher.fetchFunc = func() (external.CurrenciesDataMap, error) {
		return map[string]external.CurrencyData{
			"USD": {Symbol: "USD", Price: 1.0},
		}, nil
	}

	assert.Equal(t, 0.0, ct.data["USD"].Price)
	time.Sleep(ct.refreshInterval * 5)
	assert.Equal(t, 1.0, ct.data["USD"].Price)
	ct.stopRefreshChan <- struct{}{}
}
