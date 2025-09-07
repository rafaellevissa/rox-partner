package api_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/rafaellevissa/rox-partner/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestGetTradeSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success with start date", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		router := gin.New()
		router.GET("/summary", api.GetTradeSummary(db))

		ticker := "PETR4"
		start := "2025-09-01"

		mock.ExpectQuery(`SELECT MAX\(trade_price\).*`).
			WithArgs(ticker, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"max"}).AddRow(20.5))

		mock.ExpectQuery(`SELECT MAX\(daily_volume\).*`).
			WithArgs(ticker, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"max"}).AddRow(150000))

		req := httptest.NewRequest("GET", "/summary?ticker="+ticker+"&data_inicio="+start, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		expected := `{"ticker":"PETR4","max_range_value":20.5,"max_daily_volume":150000}`
		assert.JSONEq(t, expected, w.Body.String())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success without start date", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		router := gin.New()
		router.GET("/summary", api.GetTradeSummary(db))

		ticker := "VALE3"

		mock.ExpectQuery(`SELECT MAX\(trade_price\).*`).
			WithArgs(ticker, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"max"}).AddRow(30.0))
		mock.ExpectQuery(`SELECT MAX\(daily_volume\).*`).
			WithArgs(ticker, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"max"}).AddRow(200000))

		req := httptest.NewRequest("GET", "/summary?ticker="+ticker, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		expected := `{"ticker":"VALE3","max_range_value":30,"max_daily_volume":200000}`
		assert.JSONEq(t, expected, w.Body.String())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("missing ticker", func(t *testing.T) {
		db, _, _ := sqlmock.New()
		defer db.Close()

		router := gin.New()
		router.GET("/summary", api.GetTradeSummary(db))

		req := httptest.NewRequest("GET", "/summary", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"ticker is required"}`, w.Body.String())
	})

	t.Run("invalid start date", func(t *testing.T) {
		db, _, _ := sqlmock.New()
		defer db.Close()

		router := gin.New()
		router.GET("/summary", api.GetTradeSummary(db))

		req := httptest.NewRequest("GET", "/summary?ticker=PETR4&data_inicio=2025/09/01", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"invalid date format, use YYYY-MM-DD"}`, w.Body.String())
	})

	t.Run("db error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		defer db.Close()

		router := gin.New()
		router.GET("/summary", api.GetTradeSummary(db))

		ticker := "PETR4"

		mock.ExpectQuery(`SELECT MAX\(trade_price\).*`).
			WithArgs(ticker, sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)

		req := httptest.NewRequest("GET", "/summary?ticker="+ticker+"&data_inicio=2025-09-01", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
