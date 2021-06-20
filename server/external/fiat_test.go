package external

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFiat_Fetch(t *testing.T) {
	f := new(FiatFetcher)
	getMock := new(getAPIMock)
	f.getAPI = getMock

	type testCase struct {
		name              string
		getFiatFuncReturn []byte
		getFiatFuncError  error
		expectedError     string
		expectedTaskMap   CurrenciesDataMap
	}

	testCases := []testCase{
		{
			name:             "unknown error",
			getFiatFuncError: fmt.Errorf("unknown"),
			expectedError:    "unknown",
		},
		{
			name:              "invalid data format",
			getFiatFuncReturn: []byte(`{"rates":{"AED":4.486906,"AFN":"958a34506, "ALL":"123.228634"}}`),
			expectedError:     `invalid data format for fiat api response: invalid character 'A' after object key:value pair`,
		},
		{
			name:              "can not parse float",
			getFiatFuncReturn: []byte(`{"rates":{"AED":4.486906,"AFN":"95.8a34506", "ALL":"123.228634"}}`),
			expectedError:     "invalid data format for fiat api response: json: cannot unmarshal string into Go struct field fiatApiData.rates of type float64",
		},
		{
			name:              "ok",
			getFiatFuncReturn: []byte(`{"rates":{"AED":4.486906,"AFN":95.834506}}`),
			expectedTaskMap:   CurrenciesDataMap{"AED": CurrencyData{Symbol: "AED", Price: 4.486906}, "AFN": CurrencyData{Symbol: "AFN", Price: 95.834506}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			getMock.getFiatFunc = func() ([]byte, error) {
				return tt.getFiatFuncReturn, tt.getFiatFuncError
			}

			gotData, gotError := f.Fetch()

			if tt.getFiatFuncError == nil && tt.expectedError == "" {
				assert.NoError(t, gotError)
			} else {
				assert.EqualError(t, gotError, tt.expectedError, "kappa")
			}
			assert.Equal(t, tt.expectedTaskMap, gotData)
		})
	}
}

func TestFiat(t *testing.T) {
	c := NewFiatFetcher()
	assert.Equal(t, c.Name, "fiat")
}
