package db_portfolio

import (
	queries "github.com/matejgrzinic/portfolio/db_portfolio/queries"
)

type PortfolioAPI struct {
	CTX
}

// B A L A N C E

func (p *PortfolioAPI) GetUserBalance(user string) (*queries.DbBalanceData, error) {
	return queries.GetUserBalance(p.DB(), user)
}

func (p *PortfolioAPI) SaveBalance(b *queries.DbBalanceData) error {
	return queries.SaveBalance(p.DB(), b)
}

// T I M E L I N E

func (p *PortfolioAPI) GetUserTimeline(user string, timeframe string) (*[]queries.DbTimelineData, error) {
	return queries.GetUserTimelineQuery(p.DB(), user, timeframe)
}

// U S E R S
func (p *PortfolioAPI) GetAllUsers() (*[]queries.DbUserData, error) {
	return queries.GetAllUsers(p.DB())
}
