// Package vmwriter escreve comandos da máquina virtual do nand2tetris
// (capítulos 7 e 8 do livro), um comando por linha.
//
// A VM é uma máquina de PILHA: toda operação tira seus operandos do topo da
// pilha e coloca o resultado de volta nele. Os dados ficam em SEGMENTOS de
// memória (local, argument, this, ...), acessados com push/pop.
package vmwriter

import (
	"fmt"
	"strings"
)

// Segment é um segmento de memória da VM.
type Segment int

const (
	Constant Segment = iota // constantes 0..32767 (só push)
	Argument                // argumentos da subrotina atual
	Local                   // variáveis locais (var) da subrotina atual
	Static                  // variáveis static da classe
	This                    // campos (field) do objeto atual; base em pointer 0
	That                    // elementos de array; base em pointer 1
	Pointer                 // pointer 0 = THIS, pointer 1 = THAT
	Temp                    // 8 registradores temporários (temp 0..7)
)

// O nome do Constant no código VM é "constant", não "const".
var segmentNames = [...]string{"constant", "argument", "local", "static", "this", "that", "pointer", "temp"}

// String devolve o nome do segmento como aparece no código VM.
func (s Segment) String() string { return segmentNames[s] }

// Command é um comando aritmético/lógico da VM.
type Command int

const (
	Add Command = iota // x + y
	Sub                // x - y
	Neg                // -y   (unário)
	Eq                 // x == y  -> true (-1) ou false (0)
	Gt                 // x > y
	Lt                 // x < y
	And                // x & y   (bit a bit)
	Or                 // x | y   (bit a bit)
	Not                // ~y      (unário, bit a bit)
)

var commandNames = [...]string{"add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not"}

// String devolve o nome do comando como aparece no código VM.
func (c Command) String() string { return commandNames[c] }

// VMWriter acumula o código VM gerado para uma classe.
type VMWriter struct {
	out strings.Builder
}

// New cria um VMWriter vazio.
func New() *VMWriter { return &VMWriter{} }

// line escreve um comando e a quebra de linha.
func (w *VMWriter) line(format string, args ...any) {
	fmt.Fprintf(&w.out, format, args...)
	w.out.WriteByte('\n')
}

// WritePush escreve "push <segmento> <índice>": copia o valor para o topo da pilha.
func (w *VMWriter) WritePush(seg Segment, index int) { w.line("push %s %d", seg, index) }

// WritePop escreve "pop <segmento> <índice>": tira o topo da pilha e guarda no segmento.
func (w *VMWriter) WritePop(seg Segment, index int) { w.line("pop %s %d", seg, index) }

// WriteArithmetic escreve um comando aritmético/lógico ("add", "not", ...).
func (w *VMWriter) WriteArithmetic(cmd Command) { w.line("%s", cmd) }

// WriteLabel escreve "label <rótulo>": marca um ponto de desvio.
func (w *VMWriter) WriteLabel(label string) { w.line("label %s", label) }

// WriteGoto escreve "goto <rótulo>": desvio incondicional.
func (w *VMWriter) WriteGoto(label string) { w.line("goto %s", label) }

// WriteIf escreve "if-goto <rótulo>": tira o topo da pilha e desvia se for diferente de 0.
func (w *VMWriter) WriteIf(label string) { w.line("if-goto %s", label) }

// WriteCall escreve "call <nome> <nArgs>": chama uma função que já tem nArgs argumentos na pilha.
func (w *VMWriter) WriteCall(name string, nArgs int) { w.line("call %s %d", name, nArgs) }

// WriteFunction escreve "function <nome> <nLocals>": início de uma função com nLocals variáveis locais.
func (w *VMWriter) WriteFunction(name string, nLocals int) { w.line("function %s %d", name, nLocals) }

// WriteReturn escreve "return": devolve o topo da pilha a quem chamou.
func (w *VMWriter) WriteReturn() { w.line("return") }

// String devolve todo o código VM escrito até agora.
func (w *VMWriter) String() string { return w.out.String() }
