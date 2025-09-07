package parser

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rafaellevissa/rox-partner/internal/domain"
)

func ParseCSV(path string) ([]domain.Trade, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	r.Comma = ';'
	r.FieldsPerRecord = -1

	_, err = r.Read()
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	var trades []domain.Trade
	for {
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if len(rec) < 9 {
			continue
		}

		tradeDateStr := rec[8]
		if tradeDateStr == "" {
			tradeDateStr = rec[0]
		}
		tradeDate, err := time.Parse("2006-01-02", tradeDateStr)
		if err != nil {
			tradeDate, _ = time.Parse("02/01/2006", tradeDateStr)
		}

		priceStr := rec[3]
		priceStr = strings.Replace(priceStr, ".", "", -1)
		priceStr = strings.Replace(priceStr, ",", ".", 1)
		price, _ := strconv.ParseFloat(priceStr, 64)

		quantity, _ := strconv.Atoi(rec[4])

		closing := rec[5]
		closingFormatted := formatClosingTime(closing)

		trades = append(trades, domain.Trade{
			TradeDate:      tradeDate,
			InstrumentCode: rec[1],
			TradePrice:     price,
			TradeQuantity:  quantity,
			ClosingTime:    closingFormatted,
		})
	}

	return trades, nil
}

func formatClosingTime(s string) string {
	if len(s) < 9 {
		s = strings.Repeat("0", 9-len(s)) + s
	}
	if len(s) >= 9 {
		h := s[0:2]
		m := s[2:4]
		s2 := s[4:6]
		ms := s[6:9]
		return h + ":" + m + ":" + s2 + "." + ms
	}
	return "00:00:00.000"
}
