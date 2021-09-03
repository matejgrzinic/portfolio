package portfolio

import (
	"fmt"
	"testing"

	"github.com/matejgrzinic/portfolio/db"
	"github.com/stretchr/testify/assert"
)

func TestAllUserTransactions(t *testing.T) {
	mockedDB := db.NewMockedDB()

	mctx := &mockedCTX{db: mockedDB}

	p := NewPortfolio(mctx)
	user := &User{Name: "unittest"}

	t.Run("unknown error", func(t *testing.T) {
		mockedDB.QueryRowsFunc = func(result interface{}, rowFunc func() error) error { return fmt.Errorf("unittest") }

		_, err := p.AllUserTransactions(user)
		assert.EqualError(t, err, "unittest")
	})

	t.Run("OK", func(t *testing.T) {
		mockedDB.QueryRowsFunc = func(result interface{}, rowFunc func() error) error {
			for i := 0; i < 3; i++ {
				result.(*Transaction).Time = int64(i)
				rowFunc()
			}
			return nil
		}

		data, err := p.AllUserTransactions(user)
		assert.Len(t, data, 3)
		assert.Equal(t, int64(0), data[0].Time)
		assert.Equal(t, int64(1), data[1].Time)
		assert.Equal(t, int64(2), data[2].Time)
		assert.NoError(t, err)
	})
}

func TestInsertTransaction(t *testing.T) {
	mockedDB := db.NewMockedDB()

	mctx := &mockedCTX{db: mockedDB}

	p := NewPortfolio(mctx)

	t.Run("returns error", func(t *testing.T) {
		mockedDB.InsertOneFunc = func(data interface{}) error { return fmt.Errorf("unittest") }

		err := p.InsertTransaction(nil)
		assert.EqualError(t, err, "unittest")
	})

	inTransaction := &Transaction{
		Note:  "unittest",
		Gains: []TransactionCurrency{{Symbol: "BTC", CurrencyType: "cryptocurrency", Amount: 1.12}},
		Time:  1,
		User:  "unittest",
	}

	t.Run("OK", func(t *testing.T) {
		mockedDB.InsertOneFunc = func(data interface{}) error {
			tt, ok := data.(*Transaction)
			assert.True(t, ok)
			assert.Equal(t, inTransaction, tt)
			return nil
		}

		err := p.InsertTransaction(inTransaction)
		assert.NoError(t, err)
	})
}
