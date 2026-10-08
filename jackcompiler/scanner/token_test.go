package scanner

import (
	"strings"
	"testing"
)

// TestKeywords confere a tabela keywords contra uma lista escrita À MÃO.
// Isso pega três erros: palavra faltando, palavra sobrando/errada ("idf")
// e palavra ligada ao tipo errado (while -> DO).
func TestKeywords(t *testing.T) {
	esperado := map[string]TokenType{
		"class": CLASS, "constructor": CONSTRUCTOR, "function": FUNCTION,
		"method": METHOD, "field": FIELD, "static": STATIC, "var": VAR,
		"int": INT, "char": CHAR, "boolean": BOOLEAN, "void": VOID,
		"true": TRUE, "false": FALSE, "null": NULL, "this": THIS,
		"let": LET, "do": DO, "if": IF, "else": ELSE, "while": WHILE,
		"return": RETURN,
	}
	if len(esperado) != 21 {
		t.Fatalf("a lista do teste deveria ter 21 palavras, tem %d", len(esperado))
	}
	if len(keywords) != len(esperado) {
		t.Errorf("keywords tem %d entradas, esperado %d", len(keywords), len(esperado))
	}
	for palavra, tipo := range esperado {
		got, ok := keywords[palavra] // comma-ok: distingue "ausente" de "valor zero"
		if !ok {
			t.Errorf("keyword %q ausente da tabela", palavra)
			continue
		}
		if got != tipo {
			t.Errorf("keywords[%q] = %v, esperado %v", palavra, got, tipo)
		}
	}
}

// TestSymbols faz o mesmo para os 19 símbolos (chave rune).
func TestSymbols(t *testing.T) {
	esperado := map[rune]TokenType{
		'(': LPAREN, ')': RPAREN, '{': LBRACE, '}': RBRACE,
		'[': LBRACKET, ']': RBRACKET, ',': COMMA, ';': SEMICOLON,
		'.': DOT, '+': PLUS, '-': MINUS, '*': ASTERISK, '/': SLASH,
		'&': AND, '|': OR, '~': NOT, '<': LT, '>': GT, '=': EQ,
	}
	if len(esperado) != 19 {
		t.Fatalf("a lista do teste deveria ter 19 símbolos, tem %d", len(esperado))
	}
	if len(symbols) != len(esperado) {
		t.Errorf("symbols tem %d entradas, esperado %d", len(symbols), len(esperado))
	}
	for r, tipo := range esperado {
		got, ok := symbols[r]
		if !ok {
			t.Errorf("símbolo %q ausente da tabela", r)
			continue
		}
		if got != tipo {
			t.Errorf("symbols[%q] = %v, esperado %v", r, got, tipo)
		}
	}
}

// TestLookupIdent é um teste "table-driven": uma tabela de casos e um laço.
func TestLookupIdent(t *testing.T) {
	casos := []struct {
		palavra string
		want    TokenType
	}{
		{"while", WHILE},
		{"class", CLASS},
		{"return", RETURN},
		{"x", IDENT},
		{"contador", IDENT},
		{"whilex", IDENT}, // maximal munch: só a palavra INTEIRA conta
		{"While", IDENT},  // Jack diferencia maiúsculas de minúsculas
		{"_", IDENT},
	}
	for _, c := range casos {
		if got := lookupIdent(c.palavra); got != c.want {
			t.Errorf("lookupIdent(%q) = %v, esperado %v", c.palavra, got, c.want)
		}
	}
}

// TestCategory confere a ponte "tipo fino -> tag do XML".
func TestCategory(t *testing.T) {
	for palavra, tipo := range keywords {
		if got := tipo.Category(); got != "keyword" {
			t.Errorf("%q (%v): categoria %q, esperado keyword", palavra, tipo, got)
		}
	}
	for r, tipo := range symbols {
		if got := tipo.Category(); got != "symbol" {
			t.Errorf("%q (%v): categoria %q, esperado symbol", r, tipo, got)
		}
	}
	outros := map[TokenType]string{
		NUMBER:  "integerConstant",
		STRING:  "stringConstant",
		IDENT:   "identifier",
		ILLEGAL: "", // não aparece no XML
		EOF:     "", // não aparece no XML
	}
	for tipo, esperado := range outros {
		if got := tipo.Category(); got != esperado {
			t.Errorf("%v: categoria %q, esperado %q", tipo, got, esperado)
		}
	}
}

// TestCategoryCobreTodosOsTipos garante que NENHUM tipo "de verdade"
func TestCategoryCobreTodosOsTipos(t *testing.T) {
	for tipo := ILLEGAL + 1; tipo < EOF; tipo++ {
		if tipo.Category() == "" {
			t.Errorf("tipo %v (%d) está sem categoria", tipo, int(tipo))
		}
	}
}

// TestString confere os nomes usados na depuração.
func TestString(t *testing.T) {
	vistos := map[string]TokenType{}
	for tipo := ILLEGAL; tipo <= EOF; tipo++ {
		nome := tipo.String()
		if nome == "UNKNOWN" {
			t.Errorf("tipo %d está sem nome em tokenNames", int(tipo))
			continue
		}
		if outro, repetido := vistos[nome]; repetido {
			t.Errorf("nome %q usado por %d e por %d", nome, int(outro), int(tipo))
		}
		vistos[nome] = tipo
	}
	// Para keywords o nome é a própria palavra em maiúsculas: pega tokenNames
	// com o nome trocado (por exemplo, "IF" escrito como "ELSE").
	for palavra, tipo := range keywords {
		if got, want := tipo.String(), strings.ToUpper(palavra); got != want {
			t.Errorf("keyword %q: String() = %q, esperado %q", palavra, got, want)
		}
	}
	// Valores fora do intervalo não podem causar panic.
	for _, invalido := range []TokenType{-1, 999} {
		if got := invalido.String(); got != "UNKNOWN" {
			t.Errorf("TokenType(%d).String() = %q, esperado UNKNOWN", int(invalido), got)
		}
	}
}

// TestValorZeroEhIllegal documenta a decisão de design: um Token{} vazio
func TestValorZeroEhIllegal(t *testing.T) {
	var tok Token // todo campo começa com o valor zero do seu tipo
	if tok.Type != ILLEGAL {
		t.Errorf("Token{}.Type = %v, esperado ILLEGAL", tok.Type)
	}
	if ILLEGAL != 0 {
		t.Errorf("ILLEGAL deveria ser o valor 0, é %d", int(ILLEGAL))
	}
}