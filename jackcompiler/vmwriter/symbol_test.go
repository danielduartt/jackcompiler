package vmwriter

import (
	"testing"

	"github.com/danielduartt/jackcompiler/symtab"
)

func TestVariableAccess(t *testing.T) {
	for _, tc := range []struct {
		kind    symtab.Kind
		segment string
	}{
		{symtab.Static, "static"},
		{symtab.Field, "this"},
		{symtab.Arg, "argument"},
		{symtab.Var, "local"},
	} {
		t.Run(tc.segment, func(t *testing.T) {
			table := symtab.New()
			table.Define("x", "int", tc.kind)
			y := table.Define("y", "int", tc.kind)
			w := New()
			w.WritePushSymbol(y)
			w.WritePopSymbol(y)
			want := "push " + tc.segment + " 1\npop " + tc.segment + " 1\n"
			if got := w.String(); got != want {
				t.Fatalf("esperado %q, obtido %q", want, got)
			}
		})
	}
}
