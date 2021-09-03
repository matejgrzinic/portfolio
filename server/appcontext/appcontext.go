package appcontext

import (
	"github.com/matejgrzinic/portfolio/currencies"
	"github.com/matejgrzinic/portfolio/db"
	"github.com/matejgrzinic/portfolio/portfolio"
)

type AppContext struct {
	db db.API
	c  currencies.API
	p  portfolio.API
}

func (a *AppContext) DB() db.API {
	return a.db
}

func (a *AppContext) Currencies() currencies.API {
	return a.c
}

func (a *AppContext) Portfolio() portfolio.API {
	return a.p
}

func SetupAppContext() *AppContext {
	ctx := new(AppContext)

	db := db.NewDbAccess()
	ctx.db = db

	c := currencies.NewCurrencies(ctx)
	ctx.c = c

	p := portfolio.NewPortfolio(ctx)
	ctx.p = p

	return ctx
}
