package vmwriter

import (
	"testing"

	"github.com/danielduartt/jackcompiler/scanner"
)

func TestWriteBinaryOp(t *testing.T) {
	for _, tc := range []struct {
		op   scanner.TokenType
		want string
	}{
		{scanner.PLUS, "add"}, {scanner.MINUS, "sub"},
		{scanner.AND, "and"}, {scanner.OR, "or"},
		{scanner.LT, "lt"}, {scanner.GT, "gt"}, {scanner.EQ, "eq"},
		{scanner.ASTERISK, "call Math.multiply 2"}, {scanner.SLASH, "call Math.divide 2"},
	} {
		w := New()
		if err := w.WriteBinaryOp(tc.op); err != nil || w.String() != tc.want+"\n" {
			t.Errorf("operador %s: saída %q, esperado %q, erro %v", tc.op.Lexeme(), w.String(), tc.want, err)
		}
	}
	w := New()
	if err := w.WriteBinaryOp(scanner.NOT); err == nil || w.String() != "" {
		t.Error("operador unário deveria ser rejeitado sem emitir código")
	}
}
