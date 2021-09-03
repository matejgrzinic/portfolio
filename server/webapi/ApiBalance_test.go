package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matejgrzinic/portfolio/appcontext"
	"github.com/matejgrzinic/portfolio/portfolio"
	"github.com/stretchr/testify/assert"
)

func TestApiBalance(t *testing.T) {
	ctx := &appcontext.AppContextMock{}

	mockedPortfolio := portfolio.NewPortfolioMock()
	ctx.P = mockedPortfolio

	t.Run("no user in context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost.com/unittest", nil)
		w := httptest.NewRecorder()
		ApiBalance(ctx)(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("user in context is nil", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost.com/unittest", nil)
		req = req.WithContext(context.WithValue(req.Context(), "USER", nil))
		w := httptest.NewRecorder()
		ApiBalance(ctx)(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("UserBalance returns error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost.com/unittest", nil)
		user := &portfolio.User{Name: "unittest"}
		req = req.WithContext(context.WithValue(req.Context(), "USER", user))
		w := httptest.NewRecorder()

		mockedPortfolio.UserBalanceFunc = func(user *portfolio.User) ([]portfolio.CurrencyBalance, error) {
			assert.Equal(t, "unittest", user.Name)
			return nil, fmt.Errorf("unittest")
		}

		ApiBalance(ctx)(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("OK", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost.com/unittest", nil)
		user := &portfolio.User{Name: "unittest"}
		req = req.WithContext(context.WithValue(req.Context(), "USER", user))
		w := httptest.NewRecorder()

		mockedPortfolio.UserBalanceFunc = func(user *portfolio.User) ([]portfolio.CurrencyBalance, error) {
			assert.Equal(t, "unittest", user.Name)
			return []portfolio.CurrencyBalance{
				{CurrencyType: "1", Symbol: "1", Amount: 1.0, Price: 1.0, Value: 1.0},
				{CurrencyType: "2", Symbol: "2", Amount: 2.0, Price: 2.0, Value: 2.0},
			}, nil
		}

		ApiBalance(ctx)(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var respData struct {
			Status string                      `json:"status"`
			Data   []portfolio.CurrencyBalance `json:"data"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &respData)
		assert.NoError(t, err)
		assert.Len(t, respData.Data, 2)

		for i := 1; i <= 2; i++ {
			assert.Equal(t, fmt.Sprintf("%d", i), respData.Data[i-1].CurrencyType)
			assert.Equal(t, fmt.Sprintf("%d", i), respData.Data[i-1].Symbol)
			assert.Equal(t, float64(i), respData.Data[i-1].Amount)
			assert.Equal(t, float64(i), respData.Data[i-1].Price)
			assert.Equal(t, float64(i), respData.Data[i-1].Value)
		}
	})
}
