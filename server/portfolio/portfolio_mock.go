package portfolio

type PortfolioMock struct {
	UserBalanceFunc  func(user *User) ([]CurrencyBalance, error)
	UserTimelineFunc func(user *User, timeframe string) ([]TimelineData, error)

	InsertTransactionFunc func(t *Transaction) error
}

func NewPortfolioMock() *PortfolioMock {
	return new(PortfolioMock)
}

// B A L A N C E

func (dm *PortfolioMock) UserBalance(user *User) ([]CurrencyBalance, error) {
	return dm.UserBalanceFunc(user)
}

// T I M E L I N E

func (dm *PortfolioMock) UserTimeline(user *User, timeframe string) ([]TimelineData, error) {
	return dm.UserTimelineFunc(user, timeframe)
}

func (dm *PortfolioMock) InsertTransaction(t *Transaction) error {
	return dm.InsertTransactionFunc(t)
}
