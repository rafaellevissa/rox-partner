package parser

import (
	"os"
	"testing"
)

func TestParseCSV(t *testing.T) {
	content := `DataReferencia;CodigoInstrumento;AcaoAtualizacao;PrecoNegocio;QuantidadeNegociada;HoraFechamento;CodigoIdentificadorNegocio;TipoSessaoPregao;DataNegocio;CodigoParticipanteComprador;CodigoParticipanteVendedor
2025-09-05;SOLU25;0;207,100;1;090000004;10;1;2025-09-05;120;4090
`

	file := "test.csv"
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file)

	trades, err := ParseCSV(file)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(trades) != 1 {
		t.Fatalf("expected 1 trade, got %d", len(trades))
	}

	if trades[0].InstrumentCode != "SOLU25" {
		t.Errorf("unexpected instrument code: %s", trades[0].InstrumentCode)
	}
}
