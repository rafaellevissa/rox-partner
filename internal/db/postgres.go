package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/rafaellevissa/rox-partner/internal/domain"
)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func EnsureSchema(dbConn *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS trades (
id BIGSERIAL PRIMARY KEY,
trade_date DATE NOT NULL,
instrument_code VARCHAR(50) NOT NULL,
trade_price NUMERIC(18,6) NOT NULL,
trade_quantity INTEGER NOT NULL,
closing_time TIME(3) NOT NULL,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX IF NOT EXISTS idx_trades_instrument_date ON trades (instrument_code, trade_date);
`
	_, err := dbConn.Exec(schema)
	return err
}

func InsertTradesBulk(dbConn *sql.DB, trades []domain.Trade) error {
	if len(trades) == 0 {
		return nil
	}

	tx, err := dbConn.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("trades", "trade_date", "instrument_code", "trade_price", "trade_quantity", "closing_time"))
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, t := range trades {
		_, err = stmt.Exec(t.TradeDate.Format("2006-01-02"), t.InstrumentCode, fmt.Sprintf("%.6f", t.TradePrice), t.TradeQuantity, t.ClosingTime)
		if err != nil {
			stmt.Close()
			tx.Rollback()
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		stmt.Close()
		tx.Rollback()
		return err
	}
	if err := stmt.Close(); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func QueryTrades(dbConn *sql.DB, ticker string, startDate string, limit int, offset int) ([]domain.Trade, error) {
	query := strings.Builder{}
	query.WriteString("SELECT trade_date, instrument_code, trade_price, trade_quantity, closing_time FROM trades WHERE instrument_code = $1")
	args := []interface{}{ticker}
	argIdx := 2
	if startDate != "" {
		query.WriteString(fmt.Sprintf(" AND trade_date >= $%d", argIdx))
		args = append(args, startDate)
		argIdx++
	}
	if limit > 0 {
		query.WriteString(fmt.Sprintf(" ORDER BY trade_date, closing_time LIMIT %d OFFSET %d", limit, offset))
	}

	rows, err := dbConn.Query(query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.Trade
	for rows.Next() {
		var t domain.Trade
		var tradeDate string
		var closingTime string
		if err := rows.Scan(&tradeDate, &t.InstrumentCode, &t.TradePrice, &t.TradeQuantity, &closingTime); err != nil {
			return nil, err
		}
		tradeDate = tradeDate
		parsedDate, _ := time.Parse("2006-01-02", tradeDate)
		t.TradeDate = parsedDate
		t.ClosingTime = closingTime
		res = append(res, t)
	}

	return res, nil
}
