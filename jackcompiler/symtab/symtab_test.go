package symtab

import "testing"

// Cada Kind tem seu próprio contador: o static declarado depois de dois fields
// recebe índice 0, porque vai para outro segmento.
func TestDefineIndexesPerKind(t *testing.T) {
	st := New()
	st.Define("x", "int", Field)
	st.Define("y", "int", Field)
	st.Define("count", "int", Static)
	st.Define("size", "int", Field)

	cases := []struct {
		name  string
		typ   string
		kind  Kind
		index int
	}{
		{"x", "int", Field, 0},
		{"y", "int", Field, 1},
		{"count", "int", Static, 0},
		{"size", "int", Field, 2},
	}
	for _, c := range cases {
		sym, ok := st.Resolve(c.name)
		if !ok {
			t.Fatalf("%s não encontrado", c.name)
		}
		want := Symbol{Name: c.name, Type: c.typ, Kind: c.kind, Index: c.index}
		if sym != want {
			t.Errorf("Resolve(%s) = %+v, want %+v", c.name, sym, want)
		}
	}
	if st.VarCount(Field) != 3 || st.VarCount(Static) != 1 {
		t.Errorf("VarCount: field=%d static=%d, want 3 e 1", st.VarCount(Field), st.VarCount(Static))
	}
}

// Em um método, "this" é definido primeiro como argument 0, então o primeiro
// parâmetro declarado fica em argument 1.
func TestMethodArgsStartAtOne(t *testing.T) {
	st := New()
	st.StartSubroutine()
	st.Define("this", "Square", Arg)
	dx := st.Define("dx", "int", Arg)
	if dx.Index != 1 {
		t.Errorf("dx.Index = %d, want 1", dx.Index)
	}
}

// StartSubroutine apaga só o escopo da subrotina; o da classe continua.
func TestStartSubroutineResetsOnlySubroutineScope(t *testing.T) {
	st := New()
	st.Define("x", "int", Field)

	st.StartSubroutine()
	st.Define("a", "int", Arg)
	st.Define("i", "int", Var)

	st.StartSubroutine()
	if _, ok := st.Resolve("i"); ok {
		t.Error("i deveria ter sumido depois de StartSubroutine")
	}
	if st.VarCount(Arg) != 0 || st.VarCount(Var) != 0 {
		t.Errorf("contadores não zerados: arg=%d var=%d", st.VarCount(Arg), st.VarCount(Var))
	}
	if _, ok := st.Resolve("x"); !ok {
		t.Error("o escopo da classe não pode ser apagado")
	}
	if st.VarCount(Field) != 1 {
		t.Errorf("VarCount(Field) = %d, want 1", st.VarCount(Field))
	}
	if j := st.Define("j", "int", Var); j.Index != 0 {
		t.Errorf("j.Index = %d, want 0", j.Index)
	}
}

// Uma variável local esconde um field de mesmo nome.
func TestResolveShadowing(t *testing.T) {
	st := New()
	st.Define("x", "int", Field)
	st.StartSubroutine()
	st.Define("x", "boolean", Var)

	sym, ok := st.Resolve("x")
	if !ok || sym.Kind != Var || sym.Type != "boolean" {
		t.Errorf("Resolve(x) = %+v, want a variável local", sym)
	}
}

// Nomes que não foram declarados (classes, subrotinas) não são variáveis.
func TestResolveUnknownName(t *testing.T) {
	st := New()
	if _, ok := st.Resolve("Output"); ok {
		t.Error("Output não deveria ser resolvido como variável")
	}
}

func TestKindString(t *testing.T) {
	for kind, want := range map[Kind]string{Static: "static", Field: "field", Arg: "arg", Var: "var"} {
		if got := kind.String(); got != want {
			t.Errorf("%d.String() = %q, want %q", kind, got, want)
		}
	}
}
