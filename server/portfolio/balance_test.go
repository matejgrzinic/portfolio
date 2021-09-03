package portfolio

import (
	"fmt"
	"testing"

	"github.com/matejgrzinic/portfolio/currencies"
	"github.com/matejgrzinic/portfolio/db"
	"github.com/matejgrzinic/portfolio/external"
	"github.com/stretchr/testify/assert"
)

type mockedCTX struct {
	db         db.API
	currencies currencies.API
}

func (m *mockedCTX) DB() db.API {
	return m.db
}

func (m *mockedCTX) Currencies() currencies.API {
	return m.currencies
}

func TestUserBalance(t *testing.T) {
	mockedDB := db.NewMockedDB()
	mockedCurrencies := currencies.NewMockedCurrencies()

	mctx := &mockedCTX{
		db:         mockedDB,
		currencies: mockedCurrencies,
	}

	p := NewPortfolio(mctx)
	user := &User{Name: "unittest"}

	t.Run("AllUserTransactions returns error", func(t *testing.T) {
		mockedDB.QueryRowsFunc = func(result interface{}, rowFunc func() error) error {
			return fmt.Errorf("unittest error")
		}
		data, err := p.UserBalance(user)
		assert.EqualError(t, err, "unittest error")
		assert.Nil(t, data)
	})

	mockedDB.QueryRowsFunc = func(result interface{}, rowFunc func() error) error {
		for i := 0; i < 3; i++ {
			result.(*Transaction).Gains = []TransactionCurrency{{
				CurrencyType: "unittest",
				Symbol:       fmt.Sprintf("%d", i),
				Amount:       float64(i + 1),
			}}
			result.(*Transaction).Time = int64(i)
			rowFunc()
		}
		return nil
	}

	t.Run("GetCurrencyWithChanges returns unknown error", func(t *testing.T) {
		mockedCurrencies.GetCurrencyWithChangesFunc = func(currencyType, symbol string) (*currencies.CurrencyDataWithChanges, error) {
			return nil, fmt.Errorf("unittest error")
		}
		_, err := p.UserBalance(user)
		assert.EqualError(t, err, "unittest error")
	})

	mockedCurrencies.GetCurrencyWithChangesFunc = func(currencyType, symbol string) (*currencies.CurrencyDataWithChanges, error) {
		return &currencies.CurrencyDataWithChanges{
			CurrencyData: external.CurrencyData{Symbol: symbol, Price: 1.0},
		}, nil
	}

	t.Run("OK", func(t *testing.T) {
		data, err := p.UserBalance(user)
		assert.NoError(t, err)

		assert.Len(t, data, 3)

		for i := 0; i < 3; i++ {
			data[i].Amount = float64(i + 1)
			data[i].CurrencyType = "unittest"
			data[i].Symbol = fmt.Sprintf("%d", i)
			data[i].Price = 1.0
			data[i].Value = float64(i + 1)
		}
	})
}
