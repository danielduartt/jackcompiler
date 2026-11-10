package vmwriter

import "testing"

// Cada método do VMWriter é testado isoladamente (camada 1 do enunciado).
func TestWriterCommands(t *testing.T) {
	tests := []struct {
		name  string
		write func(w *VMWriter)
		want  string
	}{
		// push/pop em todos os segmentos
		{"push constant", func(w *VMWriter) { w.WritePush(Constant, 5) }, "push constant 5\n"},
		{"push argument", func(w *VMWriter) { w.WritePush(Argument, 0) }, "push argument 0\n"},
		{"pop local", func(w *VMWriter) { w.WritePop(Local, 1) }, "pop local 1\n"},
		{"push static", func(w *VMWriter) { w.WritePush(Static, 3) }, "push static 3\n"},
		{"pop this", func(w *VMWriter) { w.WritePop(This, 2) }, "pop this 2\n"},
		{"push that", func(w *VMWriter) { w.WritePush(That, 0) }, "push that 0\n"},
		{"pop pointer", func(w *VMWriter) { w.WritePop(Pointer, 1) }, "pop pointer 1\n"},
		{"pop temp", func(w *VMWriter) { w.WritePop(Temp, 0) }, "pop temp 0\n"},

		// comandos aritméticos e lógicos
		{"add", func(w *VMWriter) { w.WriteArithmetic(Add) }, "add\n"},
		{"sub", func(w *VMWriter) { w.WriteArithmetic(Sub) }, "sub\n"},
		{"neg", func(w *VMWriter) { w.WriteArithmetic(Neg) }, "neg\n"},
		{"eq", func(w *VMWriter) { w.WriteArithmetic(Eq) }, "eq\n"},
		{"gt", func(w *VMWriter) { w.WriteArithmetic(Gt) }, "gt\n"},
		{"lt", func(w *VMWriter) { w.WriteArithmetic(Lt) }, "lt\n"},
		{"and", func(w *VMWriter) { w.WriteArithmetic(And) }, "and\n"},
		{"or", func(w *VMWriter) { w.WriteArithmetic(Or) }, "or\n"},
		{"not", func(w *VMWriter) { w.WriteArithmetic(Not) }, "not\n"},

		// controle de fluxo
		{"label", func(w *VMWriter) { w.WriteLabel("WHILE_EXP0") }, "label WHILE_EXP0\n"},
		{"goto", func(w *VMWriter) { w.WriteGoto("WHILE_EXP0") }, "goto WHILE_EXP0\n"},
		{"if-goto", func(w *VMWriter) { w.WriteIf("WHILE_END0") }, "if-goto WHILE_END0\n"},

		// funções
		{"call", func(w *VMWriter) { w.WriteCall("Math.abs", 1) }, "call Math.abs 1\n"},
		{"function", func(w *VMWriter) { w.WriteFunction("Main.main", 2) }, "function Main.main 2\n"},
		{"return", func(w *VMWriter) { w.WriteReturn() }, "return\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			tt.write(w)
			if got := w.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// Vários comandos seguidos saem um por linha, na ordem em que foram escritos.
func TestWriterSequence(t *testing.T) {
	w := New()
	w.WriteFunction("Main.main", 0)
	w.WritePush(Constant, 0)
	w.WriteReturn()

	want := "function Main.main 0\npush constant 0\nreturn\n"
	if got := w.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
