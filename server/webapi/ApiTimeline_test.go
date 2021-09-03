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

func TestApiTimeline(t *testing.T) {
	ctx := &appcontext.AppContextMock{}

	mockedPortfolio := portfolio.NewPortfolioMock()
	ctx.P = mockedPortfolio

	t.Run("no user in context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost.com/unittest", nil)
		w := httptest.NewRecorder()
		ApiTimeline(ctx)(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("user in context is nil", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost.com/unittest", nil)
		req = req.WithContext(context.WithValue(req.Context(), "USER", nil))
		w := httptest.NewRecorder()
		ApiTimeline(ctx)(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("UserTimeframeTimeline returns error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost.com/unittest", nil)
		user := &portfolio.User{Name: "unittest"}
		req = req.WithContext(context.WithValue(req.Context(), "USER", user))
		w := httptest.NewRecorder()

		mockedPortfolio.UserTimelineFunc = func(user *portfolio.User, timeframe string) ([]portfolio.TimelineData, error) {
			assert.Equal(t, "unittest", user.Name)
			return nil, fmt.Errorf("unittest")
		}

		ApiTimeline(ctx)(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("OK", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost.com/unittest", nil)
		user := &portfolio.User{Name: "unittest"}
		req = req.WithContext(context.WithValue(req.Context(), "USER", user))
		w := httptest.NewRecorder()

		mockedPortfolio.UserTimelineFunc = func(user *portfolio.User, timeframe string) ([]portfolio.TimelineData, error) {
			assert.Equal(t, "unittest", user.Name)
			return []portfolio.TimelineData{
				{Value: 1.1, Time: 1},
				{Value: 1.2, Time: 2},
				{Value: 1.3, Time: 3},
			}, nil
		}

		ApiTimeline(ctx)(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var respData struct {
			Status string                   `json:"status"`
			Data   []portfolio.TimelineData `json:"data"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &respData)
		assert.NoError(t, err)
		assert.Len(t, respData.Data, 3)

		assert.Equal(t, respData.Data[0].Value, 1.1)
		assert.Equal(t, respData.Data[0].Time, int64(1))

		assert.Equal(t, respData.Data[1].Value, 1.2)
		assert.Equal(t, respData.Data[1].Time, int64(2))

		assert.Equal(t, respData.Data[2].Value, 1.3)
		assert.Equal(t, respData.Data[2].Time, int64(3))
	})
}
