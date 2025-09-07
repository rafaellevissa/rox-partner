package domain

import "time"

type Trade struct {
	TradeDate      time.Time `json:"trade_date"`
	InstrumentCode string    `json:"instrument_code"`
	TradePrice     float64   `json:"trade_price"`
	TradeQuantity  int       `json:"trade_quantity"`
	ClosingTime    string    `json:"closing_time"`
}
