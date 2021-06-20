package db_portfolio

import (
	queries "github.com/matejgrzinic/portfolio/db_portfolio/queries"
)

type DbMock struct {
	getUserBalanceFunc func(user string) (*queries.DbBalanceData, error)
	saveBalanceFunc    func(b *queries.DbBalanceData) error

	getUserTimelineFunc func(user string, timeframe string) (*[]queries.DbTimelineData, error)

	getAllUsersFunc func() (*[]queries.DbUserData, error)
}

// B A L A N C E

func (dm *DbMock) GetUserBalance(user string) (*queries.DbBalanceData, error) {
	return dm.getUserBalanceFunc(user)
}

func (dm *DbMock) SaveBalance(b *queries.DbBalanceData) error {
	return dm.saveBalanceFunc(b)
}

// T I M E L I N E

func (dm *DbMock) GetUserTimeline(user string, timeframe string) (*[]queries.DbTimelineData, error) {
	return dm.getUserTimelineFunc(user, timeframe)
}

// U S E R S
func (dm *DbMock) GetAllUsers() (*[]queries.DbUserData, error) {
	return dm.getAllUsersFunc()
}

func NewDbMock() *DbMock {
	dm := new(DbMock)

	dm.getUserBalanceFunc = func(user string) (*queries.DbBalanceData, error) {
		return nil, nil
	}

	dm.getAllUsersFunc = func() (*[]queries.DbUserData, error) {
		return nil, nil
	}

	dm.getUserTimelineFunc = func(user, timeframe string) (*[]queries.DbTimelineData, error) {
		return nil, nil
	}

	dm.getAllUsersFunc = func() (*[]queries.DbUserData, error) {
		return nil, nil
	}

	return dm
}
