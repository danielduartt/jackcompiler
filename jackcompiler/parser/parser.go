// Package parser implementa o analisador sintático da linguagem Jack por
// DESCIDA RECURSIVA (recursive descent): uma função para cada não-terminal da
// gramática. Ele consome os tokens do pacote scanner e gera a árvore sintática
// em XML, no formato do nand2tetris.
package parser

import (
	"fmt"
	"strings"

	"github.com/danielduartt/jackcompiler/scanner"
)

// SyntaxError é um erro de sintaxe, com a linha do token que causou o problema.
type SyntaxError struct {
	Line int
	Msg  string
}

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

// --- Controle de erros: panic/recover -----------------------------------
type abort struct{ err error }

// run prepara o lookahead, executa uma regra e converte o panic em error.
func (p *Parser) run(rule func()) (xml string, err error) {
	defer func() {
		if r := recover(); r != nil {
			a, ok := r.(abort)
			if !ok {
				panic(r) // um bug de verdade (ex.: índice fora do intervalo): não esconda
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

// --- Consumo de tokens ---------------------------------------------------

// advance consome o token seguinte: ele vira cur e o scanner fornece um novo next.
func (p *Parser) advance() {
	p.cur = p.next
	tok, err := p.scan.NextToken()
	if err != nil {
		panic(abort{fmt.Errorf("erro léxico: %w", err)})
	}
	p.next = tok
}

// peek devolve o próximo token sem consumi-lo.
func (p *Parser) peek() scanner.Token { return p.next }

// nextIs diz se o próximo token é de algum dos tipos dados (sem consumir).
func (p *Parser) nextIs(types ...scanner.TokenType) bool {
	for _, t := range types {
		if p.next.Type == t {
			return true
		}
	}
	return false
}

// match consome o próximo token SE ele for do tipo t, e diz se consumiu.
// É o jeito de tratar partes opcionais da gramática: ('else' ...)?
func (p *Parser) match(t scanner.TokenType) bool {
	if p.next.Type != t {
		return false
	}
	p.consume(t)
	return true
}

// consume EXIGE que o próximo token seja de um dos tipos dados: consome e
// escreve o token no XML. Se não for, a análise para com uma mensagem clara.
func (p *Parser) consume(types ...scanner.TokenType) {
	if !p.nextIs(types...) {
		p.fail("Esperado %s encontrado %s", describeExpected(types), describeFound(p.next))
	}
	p.advance()
	p.emit(p.cur.XML())
}

func describeExpected(types []scanner.TokenType) string {
	names := make([]string, len(types))
	for i, t := range types {
		if lx := t.Lexeme(); lx != "" {
			names[i] = "'" + lx + "'"
		} else {
			names[i] = t.Category() // identifier, integerConstant, stringConstant
		}
	}
	if len(names) == 1 {
		return names[0]
	}
	return strings.Join(names[:len(names)-1], ", ") + " ou " + names[len(names)-1]
}

func describeFound(tok scanner.Token) string {
	if tok.Type == scanner.EOF {
		return "fim do arquivo"
	}
	return "'" + tok.Value + "'"
}

// --- Saída XML -----------------------------------------------------------

// emit escreve uma linha, indentada conforme a profundidade na árvore.
func (p *Parser) emit(s string) {
	p.out.WriteString(strings.Repeat("  ", p.depth))
	p.out.WriteString(s)
	p.out.WriteByte('\n')
}

// open abre um não-terminal: <nome> e desce um nível.
func (p *Parser) open(name string) {
	p.emit("<" + name + ">")
	p.depth++
}

// close fecha um não-terminal: sobe um nível e escreve </nome>.
func (p *Parser) close(name string) {
	p.depth--
	p.emit("</" + name + ">")
}
// --- Gramática: estrutura do programa -----------------------------------

// class: 'class' className '{' classVarDec* subroutineDec* '}'
func (p *Parser) parseClass() {
	p.open("class")
	p.consume(scanner.CLASS)
	p.consume(scanner.IDENT)
	p.consume(scanner.LBRACE)
	for p.nextIs(scanner.STATIC, scanner.FIELD) {
		p.parseClassVarDec()
	}
	for p.nextIs(scanner.CONSTRUCTOR, scanner.FUNCTION, scanner.METHOD) {
		p.parseSubroutineDec()
	}
	p.consume(scanner.RBRACE)
	p.close("class")
}

// classVarDec: ('static'|'field') type varName (',' varName)* ';'
func (p *Parser) parseClassVarDec() {
	p.open("classVarDec")
	p.consume(scanner.STATIC, scanner.FIELD)
	p.parseType()
	p.consume(scanner.IDENT)
	for p.match(scanner.COMMA) {
		p.consume(scanner.IDENT)
	}
	p.consume(scanner.SEMICOLON)
	p.close("classVarDec")
}

// type: 'int' | 'char' | 'boolean' | className   (não gera tag própria)
func (p *Parser) parseType() {
	p.consume(scanner.INT, scanner.CHAR, scanner.BOOLEAN, scanner.IDENT)
}

// subroutineDec: ('constructor'|'function'|'method') ('void'|type)
//
//	subroutineName '(' parameterList ')' subroutineBody
func (p *Parser) parseSubroutineDec() {
	p.open("subroutineDec")
	p.consume(scanner.CONSTRUCTOR, scanner.FUNCTION, scanner.METHOD)
	p.consume(scanner.VOID, scanner.INT, scanner.CHAR, scanner.BOOLEAN, scanner.IDENT)
	p.consume(scanner.IDENT)
	p.consume(scanner.LPAREN)
	p.parseParameterList()
	p.consume(scanner.RPAREN)
	p.parseSubroutineBody()
	p.close("subroutineDec")
}

// parameterList: ((type varName) (',' type varName)*)?
func (p *Parser) parseParameterList() {
	p.open("parameterList")
	if !p.nextIs(scanner.RPAREN) {
		p.parseType()
		p.consume(scanner.IDENT)
		for p.match(scanner.COMMA) {
			p.parseType()
			p.consume(scanner.IDENT)
		}
	}
	p.close("parameterList")
}

// subroutineBody: '{' varDec* statements '}'
func (p *Parser) parseSubroutineBody() {
	p.open("subroutineBody")
	p.consume(scanner.LBRACE)
	for p.nextIs(scanner.VAR) {
		p.parseVarDec()
	}
	p.parseStatements()
	p.consume(scanner.RBRACE)
	p.close("subroutineBody")
}

// varDec: 'var' type varName (',' varName)* ';'
func (p *Parser) parseVarDec() {
	p.open("varDec")
	p.consume(scanner.VAR)
	p.parseType()
	p.consume(scanner.IDENT)
	for p.match(scanner.COMMA) {
		p.consume(scanner.IDENT)
	}
	p.consume(scanner.SEMICOLON)
	p.close("varDec")
}

// --- Gramática: comandos --------------------------------------------------

// statements: statement*
func (p *Parser) parseStatements() {
	p.open("statements")
	for p.nextIs(scanner.LET, scanner.IF, scanner.WHILE, scanner.DO, scanner.RETURN) {
		switch p.next.Type {
		case scanner.LET:
			p.parseLet()
		case scanner.IF:
			p.parseIf()
		case scanner.WHILE:
			p.parseWhile()
		case scanner.DO:
			p.parseDo()
		case scanner.RETURN:
			p.parseReturn()
		}
	}
	p.close("statements")
}

// letStatement: 'let' varName ('[' expression ']')? '=' expression ';'
func (p *Parser) parseLet() {
	p.open("letStatement")
	p.consume(scanner.LET)
	p.consume(scanner.IDENT)
	if p.match(scanner.LBRACKET) {
		p.parseExpression()
		p.consume(scanner.RBRACKET)
	}
	p.consume(scanner.EQ)
	p.parseExpression()
	p.consume(scanner.SEMICOLON)
	p.close("letStatement")
}

// ifStatement: 'if' '(' expression ')' '{' statements '}' ('else' '{' statements '}')?
func (p *Parser) parseIf() {
	p.open("ifStatement")
	p.consume(scanner.IF)
	p.consume(scanner.LPAREN)
	p.parseExpression()
	p.consume(scanner.RPAREN)
	p.consume(scanner.LBRACE)
	p.parseStatements()
	p.consume(scanner.RBRACE)
	if p.match(scanner.ELSE) {
		p.consume(scanner.LBRACE)
		p.parseStatements()
		p.consume(scanner.RBRACE)
	}
	p.close("ifStatement")
}

// whileStatement: 'while' '(' expression ')' '{' statements '}'
func (p *Parser) parseWhile() {
	p.open("whileStatement")
	p.consume(scanner.WHILE)
	p.consume(scanner.LPAREN)
	p.parseExpression()
	p.consume(scanner.RPAREN)
	p.consume(scanner.LBRACE)
	p.parseStatements()
	p.consume(scanner.RBRACE)
	p.close("whileStatement")
}

// doStatement: 'do' subroutineCall ';'
func (p *Parser) parseDo() {
	p.open("doStatement")
	p.consume(scanner.DO)
	p.consume(scanner.IDENT)
	p.parseCallTail()
	p.consume(scanner.SEMICOLON)
	p.close("doStatement")
}

// returnStatement: 'return' expression? ';'
func (p *Parser) parseReturn() {
	p.open("returnStatement")
	p.consume(scanner.RETURN)
	if !p.nextIs(scanner.SEMICOLON) {
		p.parseExpression()
	}
	p.consume(scanner.SEMICOLON)
	p.close("returnStatement")
}
