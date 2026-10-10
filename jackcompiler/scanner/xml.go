package scanner

import (
	"fmt"
	"io"
	"strings"
)

// strings.NewReplacer faz a troca em uma única passada, então o "&" de
// "&lt;" gerado não é escapado de novo.
var xmlEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\"", "&quot;")

// XML formata o token no padrão do nand2tetris: <categoria> valor </categoria>.
func (t Token) XML() string {
	cat := t.Type.Category()
	return fmt.Sprintf("<%s> %s </%s>", cat, xmlEscaper.Replace(t.Value), cat)
}

// WriteXML escreve a lista de tokens no formato dos arquivos *T.xml oficiais.
// io.Writer aceita arquivo, buffer, stdout... o que implementar Write.
func WriteXML(w io.Writer, tokens []Token) error {
	if _, err := fmt.Fprintln(w, "<tokens>"); err != nil {
		return err
	}
	for _, tok := range tokens {
		if _, err := fmt.Fprintln(w, tok.XML()); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(w, "</tokens>")
	return err
}
