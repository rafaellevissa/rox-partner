package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TradeSummary struct {
	Ticker         string  `json:"ticker"`
	MaxRangeValue  float64 `json:"max_range_value"`
	MaxDailyVolume int     `json:"max_daily_volume"`
}

func GetTradeSummary(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ticker := c.Query("ticker")
		startDateStr := c.Query("data_inicio")

		if ticker == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ticker is required"})
			return
		}

		var startDate time.Time
		var err error

		if startDateStr != "" {
			startDate, err = time.Parse("2006-01-02", startDateStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
				return
			}
		} else {
			end := time.Now().AddDate(0, 0, -1)
			startDate = end.AddDate(0, 0, -6)
		}

		var maxRange float64
		err = db.QueryRow(`
			SELECT MAX(trade_price)
			FROM trades
			WHERE instrument_code = $1
			  AND trade_date >= $2
		`, ticker, startDate).Scan(&maxRange)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("max price query failed: %v", err)})
			return
		}

		var maxDailyVolume int
		err = db.QueryRow(`
			SELECT MAX(daily_volume) FROM (
				SELECT SUM(trade_quantity) AS daily_volume
				FROM trades
				WHERE instrument_code = $1
				  AND trade_date >= $2
				GROUP BY trade_date
			) AS daily_totals
		`, ticker, startDate).Scan(&maxDailyVolume)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("max daily volume query failed: %v", err)})
			return
		}

		c.JSON(http.StatusOK, TradeSummary{
			Ticker:         ticker,
			MaxRangeValue:  maxRange,
			MaxDailyVolume: maxDailyVolume,
		})
	}
}
