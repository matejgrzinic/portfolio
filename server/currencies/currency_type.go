package currencies

import (
	"fmt"
	"log"
	"time"

	"github.com/matejgrzinic/portfolio/external"
)

type currencyType struct {
	data            external.CurrenciesDataMap
	refreshInterval time.Duration
	stopRefreshChan chan struct{} // test only
	fetcher         external.Fetcher
}

func newCurrecyType(fetcher external.Fetcher, refreshInterval time.Duration) *currencyType {
	ct := new(currencyType)
	ct.data = external.CurrenciesDataMap{}
	ct.fetcher = fetcher
	ct.refreshInterval = refreshInterval
	ct.stopRefreshChan = make(chan struct{}, 1)
	ct.refreshData()
	go ct.startRefresher()
	return ct
}

func (ct *currencyType) startRefresher() {
	for {
		select {
		case <-time.After(ct.refreshInterval):
			ct.refreshData()
		case <-ct.stopRefreshChan:
			return
		}
	}
}

func (ct *currencyType) refreshData() {
	data, err := ct.fetcher.Fetch()
	if err != nil {
		log.Println("refresh data:", err)
		return
	}
	ct.data = data
}

func (ct *currencyType) Get(symbol string) (*external.CurrencyData, error) {
	var c external.CurrencyData
	var ok bool
	if c, ok = ct.data[symbol]; !ok {
		return nil, fmt.Errorf("invalid symbol: %v", symbol)
	}
	return &c, nil
}
