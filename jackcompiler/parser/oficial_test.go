package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// normalize segue a recomendação do enunciado: remove os espaços nas pontas
// de cada linha e descarta as linhas vazias (e o \r do Windows, junto).
// Assim só importa a SEQUÊNCIA de tags e tokens, não a indentação.
func normalize(text string) []string {
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

// firstDiff aponta a primeira linha diferente, que indica a regra com problema.
func firstDiff(got, want []string) string {
	n := len(got)
	if len(want) < n {
		n = len(want)
	}
	for i := 0; i < n; i++ {
		if got[i] != want[i] {
			return fmt.Sprintf("linha %d:\n  nosso:   %s\n  oficial: %s", i+1, got[i], want[i])
		}
	}
	return fmt.Sprintf("as %d primeiras linhas coincidem, mas o tamanho difere: nosso=%d, oficial=%d",
		n, len(got), len(want))
}

// TestParserXMLOficial compara a árvore gerada pelo parser com os gabaritos
// oficiais do Project 10 (X.jack -> X.xml), na ordem sugerida pelo enunciado.
func TestParserXMLOficial(t *testing.T) {
	programas := []struct {
		pasta   string
		classes []string
	}{
		{"ExpressionLessSquare", []string{"Main", "Square", "SquareGame"}},
		{"Square", []string{"Main", "Square", "SquareGame"}},
		{"ArrayTest", []string{"Main"}},
	}

	for _, prog := range programas {
		for _, classe := range prog.classes {
			t.Run(prog.pasta+"/"+classe, func(t *testing.T) {
				base := filepath.Join("..", "testdata", prog.pasta, classe)

				src, err := os.ReadFile(base + ".jack")
				if err != nil {
					t.Fatalf("fonte obrigatório indisponível: %v", err)
				}
				gabarito, err := os.ReadFile(base + ".xml")
				if err != nil {
					t.Fatalf("gabarito obrigatório indisponível: %v", err)
				}

				xml, err := New(string(src)).Parse()
				if err != nil {
					t.Fatalf("erro de análise: %v", err)
				}

				got, want := normalize(xml), normalize(string(gabarito))
				if strings.Join(got, "\n") != strings.Join(want, "\n") {
					t.Errorf("XML diferente do oficial na %s", firstDiff(got, want))
				}
			})
		}
	}
}
