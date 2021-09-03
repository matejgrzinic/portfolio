package currencies

import (
	"fmt"
	"log"
	"time"

	"github.com/matejgrzinic/portfolio/external"
	"github.com/matejgrzinic/portfolio/refresher"
)

type currencyType struct {
	data            external.CurrenciesDataMap
	refreshInterval time.Duration

	errorRefreshChan chan error
	stopRefreshChan  chan struct{} // test only

	fetcher external.Fetcher
}

func newCurrecyType(fetcher external.Fetcher, refreshInterval time.Duration) *currencyType {
	ct := new(currencyType)
	ct.data = external.CurrenciesDataMap{}
	ct.fetcher = fetcher
	ct.refreshInterval = refreshInterval

	ct.errorRefreshChan = make(chan error, 1)
	ct.stopRefreshChan = make(chan struct{}, 1)

	refresher.StartRefresher(ct.errorRefreshChan, ct.stopRefreshChan, refreshInterval, ct.refreshData)
	<-ct.errorRefreshChan

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

func (ct *currencyType) refreshData() error {
	data, err := ct.fetcher.Fetch()
	if err != nil {
		log.Println("refresh data:", err)
		return err
	}
	ct.data = data
	return nil
}

func (ct *currencyType) Get(symbol string) (*external.CurrencyData, error) {
	var c external.CurrencyData
	var ok bool
	if c, ok = ct.data[symbol]; !ok {
		return nil, fmt.Errorf("invalid symbol: %v", symbol)
	}
	return &c, nil
}
