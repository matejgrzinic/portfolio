package portfolio

import (
	"github.com/matejgrzinic/portfolio/currencies"
	"github.com/matejgrzinic/portfolio/db"
)

type CTX interface {
	DB() db.API
	Currencies() currencies.API
}

type API interface {
	UserBalance(user *User) ([]CurrencyBalance, error)
	UserTimeline(user *User, timeframe string) ([]TimelineData, error)

	InsertTransaction(t *Transaction) error
}
