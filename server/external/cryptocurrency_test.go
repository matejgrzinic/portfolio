package external

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCryptoCurrency_Fetch(t *testing.T) {
	c := new(CryptocurrencyFetcher)
	externalMock := new(getAPIMock)
	c.externalApi = externalMock

	type testCase struct {
		name                        string
		getCryptocurrencyFuncReturn []byte
		getCryptocurrencyFuncError  error
		expectedError               string
		expectedTaskMap             CurrenciesDataMap
	}

	testCases := []testCase{
		{
			name:                       "unknown error",
			getCryptocurrencyFuncError: fmt.Errorf("unknown"),
			expectedError:              "unknown",
		},
		{
			name:                        "invalid data format",
			getCryptocurrencyFuncReturn: []byte(`[{"symbol":"ETHBTC","price":"0.07093900"},"symbol":"BTCUSDT","price":"35000.56232"},{"symbol":"LTCUSDT","price":"180.32322"}]`),
			expectedError:               "invalid data format for cryptocurrency api response: invalid character ':' after array element",
		},
		{
			name:                        "can not parse float",
			getCryptocurrencyFuncReturn: []byte(`[{"symbol":"ETHBTC","price":"0.07093900"},{"symbol":"BTCUSDT","price":"abc"},{"symbol":"LTCUSDT","price":"180.32322"}]`),
			expectedError:               `can not parse float from cryptocurrency price: (symbol:BTC) (value:abc): strconv.ParseFloat: parsing "abc": invalid syntax`,
		},
		{
			name:                        "ok",
			getCryptocurrencyFuncReturn: []byte(`[{"symbol":"ETHBTC","price":"0.07093900"},{"symbol":"BTCUSDT","price":"35000.56232"},{"symbol":"LTCUSDT","price":"180.32322"}]`),
			expectedTaskMap:             CurrenciesDataMap{"BTC": CurrencyData{Symbol: "BTC", Price: 35000.56232}, "LTC": CurrencyData{Symbol: "LTC", Price: 180.32322}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			externalMock.getCryptocurrencyFunc = func() ([]byte, error) {
				return tt.getCryptocurrencyFuncReturn, tt.getCryptocurrencyFuncError
			}

			gotData, gotError := c.Fetch()

			if tt.getCryptocurrencyFuncError == nil && tt.expectedError == "" {
				assert.NoError(t, gotError)
			} else {
				assert.EqualError(t, gotError, tt.expectedError, "kappa")
			}
			assert.Equal(t, tt.expectedTaskMap, gotData)
		})
	}
}

func TestNewCryptoCurrency(t *testing.T) {
	c := NewCryptocurrencyFetcher()

	assert.Equal(t, c.Name, "cryptocurrency")
}
