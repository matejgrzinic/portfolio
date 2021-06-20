package db_portfolio

type Portfolio struct {
	API
}

func NewPortfolio(ctx CTX) *Portfolio {
	return &Portfolio{API: &PortfolioAPI{CTX: ctx}}
}
