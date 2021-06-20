package db_portfolio

import (
	"github.com/matejgrzinic/portfolio/db"
	queries "github.com/matejgrzinic/portfolio/db_portfolio/queries"
)

type CTX interface {
	DB() *db.DB
}

type API interface {
	GetUserBalance(user string) (*queries.DbBalanceData, error)
	SaveBalance(b *queries.DbBalanceData) error

	GetUserTimeline(user string, timeframe string) (*[]queries.DbTimelineData, error)

	GetAllUsers() (*[]queries.DbUserData, error)
}
