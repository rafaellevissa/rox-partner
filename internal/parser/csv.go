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

func ParseCSVStream(path string, batchSize int, processBatch func([]domain.Trade) error) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	r.Comma = ';'
	r.FieldsPerRecord = -1

	if _, err := r.Read(); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	var batch []domain.Trade

	for {
		rec, err := r.Read()
		if err != nil {
			if err == io.EOF {
				if len(batch) > 0 {
					if err := processBatch(batch); err != nil {
						return err
					}
				}
				break
			}
			return err
		}

		if len(rec) < 9 {
			continue
		}

		trade := parseRecord(rec)
		batch = append(batch, trade)

		if len(batch) >= batchSize {
			if err := processBatch(batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}

	return nil
}

func parseRecord(rec []string) domain.Trade {
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

	closingFormatted := formatClosingTime(rec[5])

	return domain.Trade{
		TradeDate:      tradeDate,
		InstrumentCode: rec[1],
		TradePrice:     price,
		TradeQuantity:  quantity,
		ClosingTime:    closingFormatted,
	}
}

func formatClosingTime(s string) string {
	if len(s) < 9 {
		s = strings.Repeat("0", 9-len(s)) + s
	}
	h := s[0:2]
	m := s[2:4]
	s2 := s[4:6]
	ms := s[6:9]
	return h + ":" + m + ":" + s2 + "." + ms
}
