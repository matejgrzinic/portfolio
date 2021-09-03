package webapi

import (
	"bytes"
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

func TestApiTransactionAdd(t *testing.T) {
	ctx := &appcontext.AppContextMock{}

	mockedPortfolio := portfolio.NewPortfolioMock()
	ctx.P = mockedPortfolio

	t.Run("no user in context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost.com/unittest", nil)
		w := httptest.NewRecorder()
		ApiTransactionAdd(ctx)(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("user in context is nil", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost.com/unittest", nil)
		req = req.WithContext(context.WithValue(req.Context(), "USER", nil))
		w := httptest.NewRecorder()
		ApiTransactionAdd(ctx)(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("no transaction in request body", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost.com/unittest", nil)
		user := &portfolio.User{Name: "unittest"}
		req = req.WithContext(context.WithValue(req.Context(), "USER", user))
		w := httptest.NewRecorder()
		ApiTransactionAdd(ctx)(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("invalid data in request body", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost.com/unittest", nil)
		user := &portfolio.User{Name: "unittest"}
		req = req.WithContext(context.WithValue(req.Context(), "USER", user))
		w := httptest.NewRecorder()
		w.Body.WriteString("invalid data")
		ApiTransactionAdd(ctx)(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	b, err := json.Marshal(portfolio.Transaction{
		Time: 1,
		User: "unittest",
		Gains: []portfolio.TransactionCurrency{{
			CurrencyType: "UNITTEST",
			Symbol:       "UNITTEST",
			Amount:       1.0,
		}},
	})
	assert.NoError(t, err)

	t.Run("Insert Transaction returns error", func(t *testing.T) {
		mockedPortfolio.InsertTransactionFunc = func(tt *portfolio.Transaction) error { return fmt.Errorf("unittest") }
		req := httptest.NewRequest("GET", "http://localhost.com/unittest", bytes.NewReader(b))
		user := &portfolio.User{Name: "unittest"}
		req = req.WithContext(context.WithValue(req.Context(), "USER", user))
		w := httptest.NewRecorder()
		assert.NoError(t, err)
		ApiTransactionAdd(ctx)(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	mockedPortfolio.InsertTransactionFunc = func(tt *portfolio.Transaction) error {
		assert.Equal(t, "unittest", tt.User)
		assert.Equal(t, int64(1), tt.Time)
		return nil
	}

	t.Run("Insert Transaction returns error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost.com/unittest", bytes.NewReader(b))
		user := &portfolio.User{Name: "unittest"}
		req = req.WithContext(context.WithValue(req.Context(), "USER", user))
		w := httptest.NewRecorder()
		assert.NoError(t, err)
		ApiTransactionAdd(ctx)(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var respData struct {
			Status string     `json:"status"`
			Data   []struct{} `json:"data"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &respData)
		assert.NoError(t, err)
		assert.Equal(t, "OK", respData.Status)
		assert.Equal(t, []struct{}{}, respData.Data)
	})
}
