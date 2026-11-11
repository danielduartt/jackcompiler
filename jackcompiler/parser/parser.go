// Package parser implementa o analisador sintático da linguagem Jack por
// DESCIDA RECURSIVA (recursive descent): uma função para cada não-terminal da
// gramática. Ele consome os tokens do pacote scanner e gera a árvore sintática
// em XML, no formato do nand2tetris.
package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/danielduartt/jackcompiler/scanner"
	"github.com/danielduartt/jackcompiler/symtab"
	"github.com/danielduartt/jackcompiler/vmwriter"
)

// SyntaxError é um erro de sintaxe, com a linha do token que causou o problema.
type SyntaxError struct {
	Line   int
	Column int
	Msg    string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("Linha %d, coluna %d: %s", e.Line, e.Column, e.Msg)
}

// Parser guarda o estado da análise: o scanner, dois tokens e a saída XML.
type Parser struct {
	symbols   *symtab.SymbolTable
	scan      *scanner.Scanner
	cur       scanner.Token // último token consumido
	next      scanner.Token // lookahead: próximo token, ainda NÃO consumido
	out       strings.Builder
	vm        *vmwriter.VMWriter
	className string
	function  string // nome VM da subrotina atual: Classe.nome
	depth     int    // nível de indentação atual do XML
}

// New cria um parser para o texto-fonte src.
func New(src string) *Parser {
	return &Parser{scan: scanner.New(src)}
}

// CompileExpression traduz uma expressão de constantes inteiras, booleanas ou null.
func CompileExpression(src string) (string, error) {
	p := New(src)
	p.symbols = symtab.New() // vazia: qualquer variável na expressão é erro
	p.vm = vmwriter.New()
	_, err := p.run(func() {
		p.parseExpression()
		if !p.nextIs(scanner.EOF) {
			p.fail("Esperado fim da expressão, encontrado %s", describeFound(p.next))
		}
	})
	if err != nil {
		return "", err
	}
	return p.vm.String(), nil
}

// CompileClass gera funções void com chamadas de classe via do e retorno vazio.
func CompileClass(src string) (string, error) {
	p := New(src)
	p.symbols = symtab.New()
	p.vm = vmwriter.New()
	if _, err := p.Parse(); err != nil {
		return "", err
	}
	return p.vm.String(), nil
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
	panic(abort{&SyntaxError{Line: p.next.Line, Column: p.next.Column, Msg: fmt.Sprintf(format, args...)}})
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
	if p.vm != nil {
		return
	}
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
	p.className = p.cur.Value
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
	kind := symtab.Static
	if p.cur.Type == scanner.FIELD {
		kind = symtab.Field
	}
	p.parseType()
	typ := p.cur.Value
	p.declaredName(typ, kind)
	for p.match(scanner.COMMA) {
		p.declaredName(typ, kind)
	}
	p.consume(scanner.SEMICOLON)
	p.close("classVarDec")
}

// type: 'int' | 'char' | 'boolean' | className   (não gera tag própria)
func (p *Parser) parseType() {
	p.consume(scanner.INT, scanner.CHAR, scanner.BOOLEAN, scanner.IDENT)
}

func (p *Parser) declaredName(typ string, kind symtab.Kind) {
	p.consume(scanner.IDENT)
	if p.symbols != nil {
		p.symbols.Define(p.cur.Value, typ, kind)
	}
}

// lookup resolve na tabela de símbolos o identificador que acabou de ser
// consumido (p.cur). Um nome não declarado é erro, apontado na posição dele.
func (p *Parser) lookup() symtab.Symbol {
	sym, ok := p.symbols.Resolve(p.cur.Value)
	if !ok {
		panic(abort{&SyntaxError{Line: p.cur.Line, Column: p.cur.Column,
			Msg: fmt.Sprintf("variável '%s' não declarada", p.cur.Value)}})
	}
	return sym
}

// subroutineDec: ('constructor'|'function'|'method') ('void'|type)
//
//	subroutineName '(' parameterList ')' subroutineBody
func (p *Parser) parseSubroutineDec() {
	p.open("subroutineDec")
	if p.vm != nil {
		p.symbols.StartSubroutine()
		p.consume(scanner.FUNCTION) // constructor e method ficam para a Etapa 2
		p.consume(scanner.VOID, scanner.INT, scanner.CHAR, scanner.BOOLEAN, scanner.IDENT)
	} else {
		p.consume(scanner.CONSTRUCTOR, scanner.FUNCTION, scanner.METHOD)
		p.consume(scanner.VOID, scanner.INT, scanner.CHAR, scanner.BOOLEAN, scanner.IDENT)
	}
	p.consume(scanner.IDENT)
	p.function = p.className + "." + p.cur.Value
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
		p.declaredName(p.cur.Value, symtab.Arg)
		for p.match(scanner.COMMA) {
			p.parseType()
			p.declaredName(p.cur.Value, symtab.Arg)
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
	if p.vm != nil {
		// "function X.f nLocals" só pode ser escrito agora: o número de
		// variáveis locais só é conhecido depois de todas as declarações var.
		p.vm.WriteFunction(p.function, p.symbols.VarCount(symtab.Var))
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
	typ := p.cur.Value
	p.declaredName(typ, symtab.Var)
	for p.match(scanner.COMMA) {
		p.declaredName(typ, symtab.Var)
	}
	p.consume(scanner.SEMICOLON)
	p.close("varDec")
}

// --- Gramática: comandos --------------------------------------------------

// statements: statement*
func (p *Parser) parseStatements() {
	p.open("statements")
	for p.nextIs(scanner.LET, scanner.IF, scanner.WHILE, scanner.DO, scanner.RETURN) {
		if p.vm != nil && p.nextIs(scanner.IF, scanner.WHILE) {
			p.fail("Geração VM ainda não implementada para %s", describeFound(p.next))
		}
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
	var target symtab.Symbol
	if p.vm != nil {
		target = p.lookup()
		if p.nextIs(scanner.LBRACKET) {
			p.fail("Geração VM de arrays ainda não implementada")
		}
	}
	if p.match(scanner.LBRACKET) {
		p.parseExpression()
		p.consume(scanner.RBRACKET)
	}
	p.consume(scanner.EQ)
	p.parseExpression()
	p.consume(scanner.SEMICOLON)
	if p.vm != nil {
		p.vm.WritePopSymbol(target) // o valor da expressão vai para a variável
	}
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
	if p.vm != nil {
		p.vm.WritePop(vmwriter.Temp, 0)
	}
	p.close("doStatement")
}

// returnStatement: 'return' expression? ';'
func (p *Parser) parseReturn() {
	p.open("returnStatement")
	p.consume(scanner.RETURN)
	if !p.nextIs(scanner.SEMICOLON) {
		p.parseExpression()
	} else if p.vm != nil {
		p.vm.WritePush(vmwriter.Constant, 0) // "return;" devolve 0 (função void)
	}
	p.consume(scanner.SEMICOLON)
	if p.vm != nil {
		p.vm.WriteReturn()
	}
	p.close("returnStatement")
}

// --- Gramática: expressões -------------------------------------------------

// isOp diz se o tipo é um dos operadores binários de Jack: + - * / & | < > =
func isOp(t scanner.TokenType) bool {
	switch t {
	case scanner.PLUS, scanner.MINUS, scanner.ASTERISK, scanner.SLASH,
		scanner.AND, scanner.OR, scanner.LT, scanner.GT, scanner.EQ:
		return true
	}
	return false
}

// expression: term (op term)*
//
// Jack NÃO tem precedência de operadores: 1+2*3 é uma lista PLANA de termos e
// operadores, avaliada da esquerda para a direita. Por isso um único laço basta.
func (p *Parser) parseExpression() {
	p.open("expression")
	p.parseTerm()
	for isOp(p.next.Type) {
		op := p.next.Type
		p.consume(op)
		p.parseTerm()
		if p.vm != nil {
			if err := p.vm.WriteBinaryOp(op); err != nil {
				panic(abort{err})
			}
		}
	}
	p.close("expression")
}

// term: integerConstant | stringConstant | keywordConstant | varName
//
//	| varName '[' expression ']' | subroutineCall
//	| '(' expression ')' | unaryOp term
//
// Quando o termo começa com um identificador, o PRÓXIMO token decide:
// '[' = acesso a array; '(' ou '.' = chamada de subrotina; outro = variável.
func (p *Parser) parseTerm() {
	p.open("term")
	if p.vm != nil && p.nextIs(scanner.STRING, scanner.THIS) {
		p.fail("Geração VM ainda não implementada para %s", describeFound(p.next))
	}
	switch p.next.Type {
	case scanner.NUMBER, scanner.STRING,
		scanner.TRUE, scanner.FALSE, scanner.NULL, scanner.THIS:
		p.consume(p.next.Type)
		if p.vm != nil {
			switch p.cur.Type {
			case scanner.NUMBER:
				n, _ := strconv.Atoi(p.cur.Value) // intervalo já validado pelo scanner
				p.vm.WritePush(vmwriter.Constant, n)
			case scanner.TRUE:
				p.vm.WritePush(vmwriter.Constant, 1)
				p.vm.WriteArithmetic(vmwriter.Neg)
			case scanner.FALSE, scanner.NULL:
				p.vm.WritePush(vmwriter.Constant, 0)
			}
		}

	case scanner.IDENT:
		p.consume(scanner.IDENT)
		switch {
		case p.nextIs(scanner.LBRACKET):
			if p.vm != nil {
				p.fail("Geração VM de arrays ainda não implementada")
			}
			p.consume(scanner.LBRACKET)
			p.parseExpression()
			p.consume(scanner.RBRACKET)
		case p.nextIs(scanner.LPAREN, scanner.DOT):
			p.parseCallTail()
		default: // variável simples: empilha seu valor
			if p.vm != nil {
				p.vm.WritePushSymbol(p.lookup())
			}
		}

	case scanner.LPAREN:
		p.consume(scanner.LPAREN)
		p.parseExpression()
		p.consume(scanner.RPAREN)

	case scanner.MINUS, scanner.NOT: // operador unário
		op := p.next.Type
		p.consume(op)
		p.parseTerm()
		if p.vm != nil {
			if op == scanner.MINUS {
				p.vm.WriteArithmetic(vmwriter.Neg)
			} else {
				p.vm.WriteArithmetic(vmwriter.Not)
			}
		}

	default:
		p.fail("Esperado um termo (número, string, identificador, '(', '-' ou '~') encontrado %s",
			describeFound(p.next))
	}
	p.close("term")
}

// parseCallTail trata o RESTANTE de uma chamada de subrotina, depois que o
// primeiro identificador já foi consumido:
//
//	subroutineName '(' expressionList ')'
//	| (className|varName) '.' subroutineName '(' expressionList ')'
func (p *Parser) parseCallTail() {
	name := p.cur.Value
	if p.vm != nil {
		if _, variable := p.symbols.Resolve(name); variable || !p.nextIs(scanner.DOT) {
			p.fail("Geração VM de chamadas de método ainda não implementada")
		}
	}
	if p.match(scanner.DOT) {
		p.consume(scanner.IDENT)
		name += "." + p.cur.Value
	}
	p.consume(scanner.LPAREN)
	count := p.parseExpressionList()
	p.consume(scanner.RPAREN)
	if p.vm != nil {
		p.vm.WriteCall(name, count)
	}
}

// expressionList: (expression (',' expression)*)?
func (p *Parser) parseExpressionList() int {
	count := 0
	p.open("expressionList")
	if !p.nextIs(scanner.RPAREN) {
		p.parseExpression()
		count++
		for p.match(scanner.COMMA) {
			p.parseExpression()
			count++
		}
	}
	p.close("expressionList")
	return count
}
