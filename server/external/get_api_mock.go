package external

type getAPIMock struct {
	getCryptocurrencyFunc func() ([]byte, error)
	getFiatFunc           func() ([]byte, error)
}

func (em *getAPIMock) getCryptocurrency() ([]byte, error) {
	return em.getCryptocurrencyFunc()
}

func (em *getAPIMock) getFiat() ([]byte, error) {
	return em.getFiatFunc()
}
