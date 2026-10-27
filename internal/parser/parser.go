// Package parser analisa a gramática Jack e produz o XML do projeto 10.
package parser

import (
	"fmt"
	"strconv"
	"strings"

	"jackcompiler/internal/lexer"
)

// Error localiza o primeiro erro; Column conta bytes a partir de 1.
type Error struct {
	Line     int
	Column   int
	Expected string
	Found    lexer.Token
}

func (e *Error) Error() string {
	return fmt.Sprintf("linha %d, coluna %d: esperado %q; encontrado %s %q",
		e.Line, e.Column, e.Expected, e.Found.Type, e.Found.Literal)
}

// Parser consome o scanner existente com um token de antecipação.
// Cada instância analisa uma única classe, com uma chamada a Parse.
type Parser struct {
	lexer   *lexer.Lexer
	current lexer.Token
	column  int
	next    lexer.Token
	nextCol int
	output  strings.Builder
	depth   int
	err     error
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.next = l.NextToken()
	p.nextCol = l.TokenColumn()
	p.advance()
	return p
}

// Parse é um atalho para integrar lexer e parser a partir do código-fonte.
func Parse(input string) (string, error) {
	return New(lexer.New(input)).Parse()
}

// Parse retorna XML apenas se a classe inteira for válida, incluindo o EOF.
func (p *Parser) Parse() (string, error) {
	p.class()
	if p.err == nil && !p.match(lexer.EOF) {
		p.fail("fim da entrada após a classe")
	}
	if p.err != nil {
		return "", p.err
	}
	return p.output.String(), nil
}

func (p *Parser) advance() {
	p.current, p.column = p.next, p.nextCol
	p.next = p.lexer.NextToken()
	p.nextCol = p.lexer.TokenColumn()
}

func (p *Parser) peek() lexer.Token { return p.next }

func (p *Parser) match(tokenType lexer.TokenType) bool {
	return p.current.Type == tokenType
}

func (p *Parser) fail(expected string) {
	if p.err == nil {
		if p.match(lexer.ILLEGAL) {
			expected += " (token léxico inválido)"
		}
		p.err = &Error{p.current.Line, p.column, expected, p.current}
	}
}

func (p *Parser) consume(tokenType lexer.TokenType) {
	if p.err != nil {
		return
	}
	if !p.match(tokenType) {
		p.fail(string(tokenType))
		return
	}
	p.writeToken()
	p.advance()
}

func (p *Parser) isKeyword(words ...string) bool {
	if !p.match(lexer.KEYWORD) {
		return false
	}
	for _, word := range words {
		if p.current.Literal == word {
			return true
		}
	}
	return false
}

func (p *Parser) keyword(words ...string) {
	if p.err != nil {
		return
	}
	if !p.isKeyword(words...) {
		p.fail(strings.Join(words, " ou "))
		return
	}
	p.consume(lexer.KEYWORD)
}

var xmlEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\"", "&quot;")

func (p *Parser) writeLine(text string) {
	p.output.WriteString(strings.Repeat("  ", p.depth))
	p.output.WriteString(text)
	p.output.WriteByte('\n')
}

func (p *Parser) writeToken() {
	tag := "symbol"
	switch p.current.Type {
	case lexer.KEYWORD:
		tag = "keyword"
	case lexer.IDENTIFIER:
		tag = "identifier"
	case lexer.INT_CONST:
		tag = "integerConstant"
	case lexer.STRING_CONST:
		tag = "stringConstant"
	}
	p.writeLine("<" + tag + "> " + xmlEscaper.Replace(p.current.Literal) + " </" + tag + ">")
}

func (p *Parser) node(tag string, body func()) {
	if p.err != nil {
		return
	}
	p.writeLine("<" + tag + ">")
	p.depth++
	body()
	p.depth--
	p.writeLine("</" + tag + ">")
}

func (p *Parser) class() {
	p.node("class", func() {
		p.keyword("class")
		p.consume(lexer.IDENTIFIER)
		p.consume(lexer.LBRACE)
		for p.err == nil && p.isKeyword("static", "field") {
			p.classVarDec()
		}
		for p.err == nil && p.isKeyword("constructor", "function", "method") {
			p.subroutineDec()
		}
		p.consume(lexer.RBRACE)
	})
}

// As alternativas type e os nomes são terminais no formato XML oficial.
func (p *Parser) typeName() {
	if p.match(lexer.IDENTIFIER) {
		p.consume(lexer.IDENTIFIER)
	} else {
		p.keyword("int", "char", "boolean")
	}
}

func (p *Parser) classVarDec() {
	p.node("classVarDec", func() {
		p.keyword("static", "field")
		p.typeName()
		p.varNames()
	})
}

func (p *Parser) varNames() {
	p.consume(lexer.IDENTIFIER)
	for p.err == nil && p.match(lexer.COMMA) {
		p.consume(lexer.COMMA)
		p.consume(lexer.IDENTIFIER)
	}
	p.consume(lexer.SEMICOLON)
}

func (p *Parser) subroutineDec() {
	p.node("subroutineDec", func() {
		p.keyword("constructor", "function", "method")
		if p.isKeyword("void") {
			p.keyword("void")
		} else {
			p.typeName()
		}
		p.consume(lexer.IDENTIFIER)
		p.consume(lexer.LPAREN)
		p.parameterList()
		p.consume(lexer.RPAREN)
		p.subroutineBody()
	})
}

func (p *Parser) parameterList() {
	p.node("parameterList", func() {
		if p.match(lexer.RPAREN) {
			return
		}
		p.typeName()
		p.consume(lexer.IDENTIFIER)
		for p.err == nil && p.match(lexer.COMMA) {
			p.consume(lexer.COMMA)
			p.typeName()
			p.consume(lexer.IDENTIFIER)
		}
	})
}

func (p *Parser) subroutineBody() {
	p.node("subroutineBody", func() {
		p.consume(lexer.LBRACE)
		for p.err == nil && p.isKeyword("var") {
			p.varDec()
		}
		p.statements()
		p.consume(lexer.RBRACE)
	})
}

func (p *Parser) varDec() {
	p.node("varDec", func() {
		p.keyword("var")
		p.typeName()
		p.varNames()
	})
}

func (p *Parser) statements() {
	p.node("statements", func() {
		for p.err == nil && p.isKeyword("let", "if", "while", "do", "return") {
			p.statement()
		}
	})
}

func (p *Parser) statement() {
	switch p.current.Literal {
	case "let":
		p.letStatement()
	case "if":
		p.ifStatement()
	case "while":
		p.whileStatement()
	case "do":
		p.doStatement()
	case "return":
		p.returnStatement()
	}
}

func (p *Parser) letStatement() {
	p.node("letStatement", func() {
		p.keyword("let")
		p.consume(lexer.IDENTIFIER)
		if p.match(lexer.LBRACKET) {
			p.consume(lexer.LBRACKET)
			p.expression()
			p.consume(lexer.RBRACKET)
		}
		p.consume(lexer.EQ)
		p.expression()
		p.consume(lexer.SEMICOLON)
	})
}

func (p *Parser) ifStatement() {
	p.node("ifStatement", func() {
		p.keyword("if")
		p.condition()
		p.statementBlock()
		if p.isKeyword("else") {
			p.keyword("else")
			p.statementBlock()
		}
	})
}

func (p *Parser) whileStatement() {
	p.node("whileStatement", func() {
		p.keyword("while")
		p.condition()
		p.statementBlock()
	})
}

func (p *Parser) condition() {
	p.consume(lexer.LPAREN)
	p.expression()
	p.consume(lexer.RPAREN)
}

func (p *Parser) statementBlock() {
	p.consume(lexer.LBRACE)
	p.statements()
	p.consume(lexer.RBRACE)
}

func (p *Parser) doStatement() {
	p.node("doStatement", func() {
		p.keyword("do")
		p.subroutineCall()
		p.consume(lexer.SEMICOLON)
	})
}

func (p *Parser) returnStatement() {
	p.node("returnStatement", func() {
		p.keyword("return")
		if !p.match(lexer.SEMICOLON) {
			p.expression()
		}
		p.consume(lexer.SEMICOLON)
	})
}

// Jack define expression = term (op term)*. Operadores binários ficam
// no mesmo nível; parênteses e unaryOp introduzem termos aninhados.
func (p *Parser) expression() {
	p.node("expression", func() {
		p.term()
		for p.err == nil && p.isOp() {
			p.op()
			p.term()
		}
	})
}

func (p *Parser) isOp() bool {
	switch p.current.Type {
	case lexer.PLUS, lexer.MINUS, lexer.ASTERISK, lexer.SLASH, lexer.AND,
		lexer.OR, lexer.LT, lexer.GT, lexer.EQ:
		return true
	}
	return false
}

func (p *Parser) op() { p.consume(p.current.Type) }

func (p *Parser) unaryOp() { p.consume(p.current.Type) }

func (p *Parser) keywordConstant() { p.keyword("true", "false", "null", "this") }

func (p *Parser) term() {
	p.node("term", func() {
		switch p.current.Type {
		case lexer.INT_CONST:
			n, err := strconv.Atoi(p.current.Literal)
			if err != nil || n > 32767 {
				p.fail("inteiro entre 0 e 32767")
				return
			}
			p.consume(lexer.INT_CONST)
		case lexer.STRING_CONST:
			p.consume(lexer.STRING_CONST)
		case lexer.KEYWORD:
			p.keywordConstant()
		case lexer.IDENTIFIER:
			if p.peek().Type == lexer.LPAREN || p.peek().Type == lexer.DOT {
				p.subroutineCall()
			} else {
				p.consume(lexer.IDENTIFIER)
				if p.match(lexer.LBRACKET) {
					p.consume(lexer.LBRACKET)
					p.expression()
					p.consume(lexer.RBRACKET)
				}
			}
		case lexer.LPAREN:
			p.consume(lexer.LPAREN)
			p.expression()
			p.consume(lexer.RPAREN)
		case lexer.MINUS, lexer.NOT:
			p.unaryOp()
			p.term()
		default:
			p.fail("termo (constante, variável, chamada, parênteses ou operador unário)")
		}
	})
}

func (p *Parser) subroutineCall() {
	p.consume(lexer.IDENTIFIER)
	if p.match(lexer.DOT) {
		p.consume(lexer.DOT)
		p.consume(lexer.IDENTIFIER)
	}
	p.consume(lexer.LPAREN)
	p.expressionList()
	p.consume(lexer.RPAREN)
}

func (p *Parser) expressionList() {
	p.node("expressionList", func() {
		if p.match(lexer.RPAREN) {
			return
		}
		p.expression()
		for p.err == nil && p.match(lexer.COMMA) {
			p.consume(lexer.COMMA)
			p.expression()
		}
	})
}
