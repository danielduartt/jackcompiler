// Package symtab implementa a tabela de símbolos do compilador Jack.
//
// O código VM não tem nomes de variáveis, só segmentos e índices
// (local 0, this 2, ...). A tabela de símbolos faz essa tradução: para cada
// identificador ela guarda o tipo, o Kind (static, field, arg ou var) e o
// índice dentro do seu segmento.
package symtab

// Kind é a categoria de uma variável. Cada Kind vira um segmento da VM:
// Static -> static, Field -> this, Arg -> argument, Var -> local.
type Kind int

const (
	Static Kind = iota // variável de classe compartilhada (static)
	Field              // campo do objeto (field)
	Arg                // parâmetro da subrotina
	Var                // variável local da subrotina (var)
)

var kindNames = [...]string{"static", "field", "arg", "var"}

// String devolve o nome do Kind, útil em mensagens e testes.
func (k Kind) String() string { return kindNames[k] }

// Symbol é tudo o que o gerador de código precisa saber sobre um identificador.
type Symbol struct {
	Name  string
	Type  string // int, char, boolean ou um nome de classe (Square, Array, ...)
	Kind  Kind
	Index int // posição dentro do segmento do Kind
}
