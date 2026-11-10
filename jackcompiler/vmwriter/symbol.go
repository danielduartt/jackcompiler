package vmwriter

import "github.com/danielduartt/jackcompiler/symtab"

// Um índice é relativo ao seu segmento. Um field e um local podem ter índice 0.
var variableSegments = [...]Segment{
	symtab.Static: Static,
	symtab.Field:  This,
	symtab.Arg:    Argument,
	symtab.Var:    Local,
}

// WritePushSymbol empilha o valor da variável resolvida pela tabela.
func (w *VMWriter) WritePushSymbol(symbol symtab.Symbol) {
	w.WritePush(variableSegments[symbol.Kind], symbol.Index)
}

// WritePopSymbol guarda o topo da pilha na variável resolvida pela tabela.
func (w *VMWriter) WritePopSymbol(symbol symtab.Symbol) {
	w.WritePop(variableSegments[symbol.Kind], symbol.Index)
}
