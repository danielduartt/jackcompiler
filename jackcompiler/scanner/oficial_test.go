package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// normalizar remove diferenças irrelevantes entre sistemas: o fim de linha do
// Windows (\r\n) e quebras/espaços sobrando no começo ou no fim do arquivo.
// O CONTEÚDO de cada linha continua sendo comparado exatamente.
func normalizar(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.TrimSpace(s)
}

// primeiraDiferenca descreve onde dois textos começam a divergir,
// para que a mensagem de falha aponte o problema em vez de só dizer "diferente".
func primeiraDiferenca(got, want string) string {
	g := strings.Split(got, "\n")
	w := strings.Split(want, "\n")
	menor := len(g)
	if len(w) < menor {
		menor = len(w)
	}
	for i := 0; i < menor; i++ {
		if g[i] != w[i] {
			return fmt.Sprintf("linha %d:\n  nosso:   %q\n  oficial: %q", i+1, g[i], w[i])
		}
	}
	return fmt.Sprintf("as %d primeiras linhas coincidem, mas o tamanho difere: nosso=%d linhas, oficial=%d linhas",
		menor, len(g), len(w))
}

// TestXMLOficial percorre ../testdata e, para cada X.jack que tenha um
// gabarito XT.xml ao lado, confere se a saída do scanner é IDÊNTICA a ele.
func TestXMLOficial(t *testing.T) {
	raiz := filepath.Join("..", "testdata")
	if _, err := os.Stat(raiz); err != nil {
		t.Skipf("pasta %s não encontrada: coloque os arquivos oficiais nela", raiz)
	}

	comparados := 0
	err := filepath.WalkDir(raiz, func(caminho string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(caminho, ".jack") {
			return nil
		}
		gabarito := strings.TrimSuffix(caminho, ".jack") + "T.xml"
		if _, err := os.Stat(gabarito); err != nil {
			return nil // este .jack não tem XML oficial ao lado: ignora
		}
		comparados++

		nome, _ := filepath.Rel(raiz, caminho)
		t.Run(filepath.ToSlash(nome), func(t *testing.T) {
			src, err := os.ReadFile(caminho)
			if err != nil {
				t.Fatal(err)
			}
			tokens, err := New(string(src)).Tokens()
			if err != nil {
				t.Fatalf("erro léxico: %v", err)
			}
			var sb strings.Builder
			if err := WriteXML(&sb, tokens); err != nil {
				t.Fatal(err)
			}
			esperado, err := os.ReadFile(gabarito)
			if err != nil {
				t.Fatal(err)
			}

			got, want := normalizar(sb.String()), normalizar(string(esperado))
			if got != want {
				t.Errorf("saída diferente da oficial (%s)\nprimeira diferença em %s",
					filepath.Base(gabarito), primeiraDiferenca(got, want))
				return
			}
			t.Logf("idêntico ao oficial: %d tokens", len(tokens))
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if comparados == 0 {
		t.Skip("nenhum par X.jack + XT.xml encontrado em ../testdata")
	}
}
