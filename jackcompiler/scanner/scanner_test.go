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
			{LET, "let", 1, 1}, {IDENT, "x", 1, 5}, {EQ, "=", 1, 7}, {NUMBER, "10", 1, 9}, {SEMICOLON, ";", 1, 11},
		}},
		{"símbolos colados", "if(x<10){}", []Token{
			{IF, "if", 1, 1}, {LPAREN, "(", 1, 3}, {IDENT, "x", 1, 4}, {LT, "<", 1, 5}, {NUMBER, "10", 1, 6},
			{RPAREN, ")", 1, 8}, {LBRACE, "{", 1, 9}, {RBRACE, "}", 1, 10},
		}},
		{"comentário de linha", "// oi\nreturn;", []Token{
			{RETURN, "return", 2, 1}, {SEMICOLON, ";", 2, 7},
		}},
		{"comentário de bloco", "/* a\nb */ do", []Token{
			{DO, "do", 2, 6},
		}},
		{"comentário de documentação", "/** doc */ var", []Token{
			{VAR, "var", 1, 12},
		}},
		{"divisão não é comentário", "a / b", []Token{
			{IDENT, "a", 1, 1}, {SLASH, "/", 1, 3}, {IDENT, "b", 1, 5},
		}},
		{"// dentro de string", "\"http://x\";", []Token{
			{STRING, "http://x", 1, 1}, {SEMICOLON, ";", 1, 11},
		}},
		{"string com espaços", "\"Hello World\"", []Token{
			{STRING, "Hello World", 1, 1},
		}},
		{"identificadores parecidos com keyword", "while_x whilex _a1", []Token{
			{IDENT, "while_x", 1, 1}, {IDENT, "whilex", 1, 9}, {IDENT, "_a1", 1, 16},
		}},
		{"maior inteiro válido", "32767", []Token{
			{NUMBER, "32767", 1, 1},
		}},
		{"linhas contadas corretamente", "let\n\nx\n", []Token{
			{LET, "let", 1, 1}, {IDENT, "x", 3, 1},
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
		{"string com quebra CR", "\"a\rb\""},
		{"string com quebra CRLF", "\"a\r\nb\""},
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
		if tok, err := s.NextToken(); err != nil || tok.Type != EOF || tok.Column != 2 {
			t.Errorf("esperava EOF, veio %v (err=%v)", tok, err)
		}
	}
}

func TestPosicaoDosErrosLexicos(t *testing.T) {
	for _, src := range []string{
		"// início\r\n\t/*é*/ @",
		"// início\r\t/*é*/ 32768",
		"// início\n\t/*é*/ \"sem fim",
		"// início\n\t/*é*/ /*aberto\nna próxima linha",
	} {
		tok, err := New(src).NextToken()
		e, ok := err.(*LexError)
		if !ok || e.Line != 2 || e.Column != 8 || !strings.Contains(e.Error(), "linha 2, coluna 8") {
			t.Fatalf("%q: esperado diagnóstico na linha 2, coluna 8; obtido %v", src, err)
		}
		if tok.Type != ILLEGAL || tok.Line != e.Line || tok.Column != e.Column {
			t.Errorf("posição do token inválido diverge do diagnóstico: %+v", tok)
		}
	}
}

func TestWriteXMLOmiteEOF(t *testing.T) {
	var sb strings.Builder
	tokens := []Token{{Type: IDENT, Value: "x", Line: 1}, {Type: EOF, Line: 1}}
	if err := WriteXML(&sb, tokens); err != nil {
		t.Fatal(err)
	}
	want := "<tokens>\n<identifier> x </identifier>\n</tokens>\n"
	if sb.String() != want {
		t.Errorf("XML com EOF: got %q, want %q", sb.String(), want)
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
