package parser

import (
	"fmt"
	"strings"

	"github.com/danielduartt/jackcompiler/scanner"
)

// SyntaxError é um erro de sintaxe, com a linha do token e a mensagem de erro.
type SyntaxError struct {
	Line int
	Msg  string
}

// Função de formatação da mensagem de erro de sintaxe.
func (e *SyntaxError) Error() string {
	return fmt.Sprintf("Linha %d: %s", e.Line, e.Msg)
}

// Parser guarda o estado da análise: o scanner, dois tokens e a saída XML.
type Parser struct {
	scan  *scanner.Scanner
	cur   scanner.Token // último token consumido
	next  scanner.Token // lookahead: próximo token, ainda NÃO consumido
	out   strings.Builder
	depth int // nível de indentação atual do XML
}

// New cria um parser para o texto-fonte src.
func New(src string) *Parser {
	return &Parser{scan: scanner.New(src)}
}

// Parse analisa uma classe Jack completa e devolve a árvore sintática em XML.
// No primeiro erro, a análise para e o erro é devolvido (sem recuperação).
func (p *Parser) Parse() (string, error) {
	return p.run(func() {
		p.parseClass()
		if p.next.Type != scanner.EOF {
			p.fail("Esperado fim do arquivo encontrado %s", describeFound(p.next))
		}
	})
}

// --- Controle de erros: panic/recover ------
// Numa descida recursiva, um erro pode surgir 10 chamadas "dentro" da gramática.
// Em vez de devolver (valor, error) em TODA função, usamos panic para abortar e
// recover num único lugar (run).

// abort é o valor lançado com panic para interromper a análise.
type abort struct{ err error }

// run prepara o lookahead, executa uma regra e converte o panic em error.
func (p *Parser) run(rule func()) (xml string, err error) {
	defer func() {
		if r := recover(); r != nil {
			a, ok := r.(abort)
			if !ok {
				panic(r)
			}
			xml, err = "", a.err
		}
	}()
	p.advance() // carrega o primeiro token em p.next
	rule()
	return p.out.String(), nil
}

// fail interrompe a análise com um erro de sintaxe na linha do próximo token.
func (p *Parser) fail(format string, args ...any) {
	panic(abort{&SyntaxError{Line: p.next.Line, Msg: fmt.Sprintf(format, args...)}})
}
