package scanner

import (
	"reflect"
	"strings"
	"testing"
)

func TestTokens(t *testing.T) {
	casos := []struct {
		nome string
		src  string
		want []Token
	}{
		{"vazio", "", nil},
		{"só espaços", " \t\r\n ", nil},
		{"let simples", "let x = 10;", []Token{
			{LET, "let", 1}, {IDENT, "x", 1}, {EQ, "=", 1}, {NUMBER, "10", 1}, {SEMICOLON, ";", 1},
		}},
		{"símbolos colados", "if(x<10){}", []Token{
			{IF, "if", 1}, {LPAREN, "(", 1}, {IDENT, "x", 1}, {LT, "<", 1}, {NUMBER, "10", 1},
			{RPAREN, ")", 1}, {LBRACE, "{", 1}, {RBRACE, "}", 1},
		}},
		{"comentário de linha", "// oi\nreturn;", []Token{
			{RETURN, "return", 2}, {SEMICOLON, ";", 2},
		}},
		{"comentário de bloco", "/* a\nb */ do", []Token{
			{DO, "do", 2},
		}},
		{"comentário de documentação", "/** doc */ var", []Token{
			{VAR, "var", 1},
		}},
		{"divisão não é comentário", "a / b", []Token{
			{IDENT, "a", 1}, {SLASH, "/", 1}, {IDENT, "b", 1},
		}},
		{"// dentro de string", "\"http://x\";", []Token{
			{STRING, "http://x", 1}, {SEMICOLON, ";", 1},
		}},
		{"string com espaços", "\"Hello World\"", []Token{
			{STRING, "Hello World", 1},
		}},
		{"identificadores parecidos com keyword", "while_x whilex _a1", []Token{
			{IDENT, "while_x", 1}, {IDENT, "whilex", 1}, {IDENT, "_a1", 1},
		}},
		{"maior inteiro válido", "32767", []Token{
			{NUMBER, "32767", 1},
		}},
		{"linhas contadas corretamente", "let\n\nx\n", []Token{
			{LET, "let", 1}, {IDENT, "x", 3},
		}},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			got, err := New(c.src).Tokens()
			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("\n got: %v\nwant: %v", got, c.want)
			}
		})
	}
}

func TestErros(t *testing.T) {
	casos := []struct{ nome, src string }{
		{"bloco não fechado", "/* nunca fecha"},
		{"string não fechada", "\"sem fim"},
		{"string quebrada em linhas", "\"a\nb\""},
		{"caractere inválido", "let x = #;"},
		{"inteiro grande demais", "32768"},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			if _, err := New(c.src).Tokens(); err == nil {
				t.Errorf("esperava erro para %q", c.src)
			}
		})
	}
}

func TestEOF(t *testing.T) {
	s := New("x")
	if tok, _ := s.NextToken(); tok.Type != IDENT {
		t.Fatalf("primeiro token: %v", tok)
	}
	for i := 0; i < 3; i++ { // chamar de novo no fim deve sempre dar EOF
		if tok, err := s.NextToken(); err != nil || tok.Type != EOF {
			t.Errorf("esperava EOF, veio %v (err=%v)", tok, err)
		}
	}
}

func TestWriteXML(t *testing.T) {
	toks, err := New("if (x < 1) { do a.b(\"a&b\"); }").Tokens()
	if err != nil {
		t.Fatal(err)
	}
	var sb strings.Builder
	if err := WriteXML(&sb, toks); err != nil {
		t.Fatal(err)
	}
	want := strings.Join([]string{
		"<tokens>",
		"<keyword> if </keyword>",
		"<symbol> ( </symbol>",
		"<identifier> x </identifier>",
		"<symbol> &lt; </symbol>",
		"<integerConstant> 1 </integerConstant>",
		"<symbol> ) </symbol>",
		"<symbol> { </symbol>",
		"<keyword> do </keyword>",
		"<identifier> a </identifier>",
		"<symbol> . </symbol>",
		"<identifier> b </identifier>",
		"<symbol> ( </symbol>",
		"<stringConstant> a&amp;b </stringConstant>",
		"<symbol> ) </symbol>",
		"<symbol> ; </symbol>",
		"<symbol> } </symbol>",
		"</tokens>",
		"",
	}, "\n")
	if sb.String() != want {
		t.Errorf("\n got:\n%s\nwant:\n%s", sb.String(), want)
	}
}
