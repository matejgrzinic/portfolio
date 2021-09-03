//+build integration

package portfolio

import (
	"fmt"
	"testing"
	"time"

	"github.com/matejgrzinic/portfolio/db"
	"github.com/stretchr/testify/assert"
)

func TestAbcD(t *testing.T) {
	ctx := &mockedCTX{db: db.NewDbAccess()}

	tt := &Transaction{
		Time: time.Now().Unix(),
		User: "Ace",
		Gains: []TransactionCurrency{
			{
				CurrencyType: "cryptocurrency",
				Symbol:       "BTC",
				Amount:       1.12,
			},
		},
		Losses: []TransactionCurrency{},
	}

	p := NewPortfolio(ctx)
	p.InsertTransaction(tt)
}

func TestAbcD2(t *testing.T) {
	ctx := &mockedCTX{db: db.NewDbAccess()}

	p := NewPortfolio(ctx)

	tts, err := p.AllUserTransactions(&User{Name: "Ace"})

	fmt.Println(tts)
	fmt.Println(err)
	assert.Fail(t, "kek")
}
