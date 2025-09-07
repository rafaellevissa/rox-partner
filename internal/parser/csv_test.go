package parser

import (
	"os"
	"testing"
	"time"

	"github.com/rafaellevissa/rox-partner/internal/domain"
)

func TestParseCSVStream(t *testing.T) {
	content := `DataReferencia;CodigoInstrumento;AcaoAtualizacao;PrecoNegocio;QuantidadeNegociada;HoraFechamento;CodigoIdentificadorNegocio;TipoSessaoPregao;DataNegocio;CodigoParticipanteComprador;CodigoParticipanteVendedor
2025-09-05;SOLU25;0;207,100;1;090000004;10;1;2025-09-05;120;4090
`

	file := "test.csv"
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file)

	batchSize := 1
	batches := [][]domain.Trade{}

	err := ParseCSVStream(file, batchSize, func(batch []domain.Trade) error {
		batches = append(batches, batch)
		return nil
	})

	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(batches) != 1 {
		t.Fatalf("expected 1 batch, got %d", len(batches))
	}

	if len(batches[0]) != 1 {
		t.Fatalf("expected 1 trade in batch, got %d", len(batches[0]))
	}

	trade := batches[0][0]
	if trade.InstrumentCode != "SOLU25" {
		t.Errorf("unexpected instrument code: %s", trade.InstrumentCode)
	}
	if trade.TradeQuantity != 1 {
		t.Errorf("unexpected quantity: %d", trade.TradeQuantity)
	}
	if trade.TradePrice != 207.100 {
		t.Errorf("unexpected price: %.3f", trade.TradePrice)
	}
	expectedDate := time.Date(2025, 9, 5, 0, 0, 0, 0, time.UTC)
	if !trade.TradeDate.Equal(expectedDate) {
		t.Errorf("unexpected trade date: %v", trade.TradeDate)
	}
	if trade.ClosingTime != "09:00:00.004" {
		t.Errorf("unexpected closing time: %s", trade.ClosingTime)
	}
}
