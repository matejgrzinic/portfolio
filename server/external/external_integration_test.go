// +build integration

package external

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getCryptocurrency(t *testing.T) {
	e := new(getImpl)

	t.Run("OK", func(t *testing.T) {
		b, err := e.getCryptocurrency()
		assert.NoError(t, err)
		assert.NotNil(t, b)
	})

	t.Run("invalid url", func(t *testing.T) {
		tmp := os.Getenv("CRYPTOCURRENCY_URL")
		os.Setenv("CRYPTOCURRENCY_URL", "x")
		defer func() {
			os.Setenv("CRYPTOCURRENCY_URL", tmp)
		}()

		b, err := e.getCryptocurrency()
		assert.EqualError(t, err, `get cryptocurrency api: Get "x": unsupported protocol scheme ""`)
		assert.Nil(t, b)
	})
}

func TestCryptocurrencyFetcher(t *testing.T) {
	cf := NewCryptocurrencyFetcher()

	data, err := cf.Fetch()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Greater(t, data["BTC"].Price, 0.0)
}

func Test_getFiat(t *testing.T) {
	e := new(getImpl)

	t.Run("OK", func(t *testing.T) {
		b, err := e.getFiat()
		assert.NoError(t, err)
		assert.NotNil(t, b)
	})

	t.Run("invalid url", func(t *testing.T) {
		tmp := os.Getenv("FIAT_URL")
		os.Setenv("FIAT_URL", "x")
		defer func() {
			os.Setenv("FIAT_URL", tmp)
		}()

		b, err := e.getFiat()
		assert.EqualError(t, err, `get fiat api: Get "x": unsupported protocol scheme ""`)
		assert.Nil(t, b)
	})
}

func TestFiatFetcher(t *testing.T) {
	ff := NewFiatFetcher()

	data, err := ff.Fetch()
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Greater(t, data["USD"].Price, 0.0)
	assert.Equal(t, data["EUR"].Price, 1.0)
}
