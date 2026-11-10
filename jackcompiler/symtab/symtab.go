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

// SymbolTable tem dois escopos:
//   - classe: static e field, válidos na classe inteira;
//   - subrotina: arg e var, válidos só na subrotina atual.
//
// Há um contador por Kind: o índice de um símbolo novo é quantos do mesmo
// Kind já existem. Assim um static declarado depois de dois fields recebe
// índice 0, porque fica em outro segmento.
type SymbolTable struct {
	class      map[string]Symbol
	subroutine map[string]Symbol
	counts     [len(kindNames)]int
}

// New cria uma tabela vazia. Use uma tabela por classe.
func New() *SymbolTable {
	return &SymbolTable{
		class:      map[string]Symbol{},
		subroutine: map[string]Symbol{},
	}
}

// StartSubroutine começa uma subrotina nova: apaga o escopo da subrotina
// anterior e zera os contadores de arg e var. O escopo da classe continua.
func (st *SymbolTable) StartSubroutine() {
	st.subroutine = map[string]Symbol{}
	st.counts[Arg] = 0
	st.counts[Var] = 0
}

// Resolve procura um nome primeiro no escopo da subrotina e depois no da
// classe, de modo que uma variável local "esconde" um field de mesmo nome.
// ok == false indica que o nome não é variável: é um nome de classe
// (Output.printInt) ou de subrotina.
func (st *SymbolTable) Resolve(name string) (sym Symbol, ok bool) {
	if sym, ok = st.subroutine[name]; ok {
		return sym, true
	}
	sym, ok = st.class[name]
	return sym, ok
}

// Define adiciona um identificador ao escopo do seu Kind e devolve o símbolo
// criado, já com o índice.
func (st *SymbolTable) Define(name, typ string, kind Kind) Symbol {
	sym := Symbol{Name: name, Type: typ, Kind: kind, Index: st.counts[kind]}
	st.counts[kind]++
	if kind == Static || kind == Field {
		st.class[name] = sym
	} else {
		st.subroutine[name] = sym
	}
	return sym
}

// VarCount devolve quantos símbolos do Kind já foram definidos no escopo atual.
// O gerador usa VarCount(Var) no "function X.f nLocals" e VarCount(Field) no
// tamanho do objeto alocado pelo construtor.
func (st *SymbolTable) VarCount(kind Kind) int { return st.counts[kind] }
