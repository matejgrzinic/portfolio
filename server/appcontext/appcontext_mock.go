package appcontext

import (
	"github.com/matejgrzinic/portfolio/currencies"
	"github.com/matejgrzinic/portfolio/db"
	"github.com/matejgrzinic/portfolio/portfolio"
)

type AppContextMock struct {
	Db db.API
	C  currencies.API
	P  portfolio.API
}

func (am *AppContextMock) DB() db.API {
	return am.Db
}

func (am *AppContextMock) Currencies() currencies.API {
	return am.C
}

func (am *AppContextMock) Portfolio() portfolio.API {
	return am.P
}
