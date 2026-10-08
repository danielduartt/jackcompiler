package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let x = 10;`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{KEYWORD, "let"},
		{IDENTIFIER, "x"},
		{EQ, "="},
		{INT_CONST, "10"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - TokenType errado. esperado=%q, obteve=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - Literal errado. esperado=%q, obteve=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestStringTokens(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		literal string
		line    int
	}{
		{"normal", `"Hello, World!"`, "Hello, World!", 1},
		{"vazia", `""`, "", 1},
		{"simbolos e comentarios", `"{} () [] + - * / & | < > = ~ ; // /* */"`, "{} () [] + - * / & | < > = ~ ; // /* */", 1},
		{"barras literais", `"\n\t"`, `\n\t`, 1},
		{"linha da abertura", "\n\n\"texto\"", "texto", 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			want := Token{Type: STRING_CONST, Literal: tt.literal, Line: tt.line}
			if got := l.NextToken(); got != want {
				t.Fatalf("token esperado=%+v, obteve=%+v", want, got)
			}
			want = Token{Type: EOF, Literal: "", Line: tt.line}
			if got := l.NextToken(); got != want {
				t.Fatalf("token seguinte esperado=%+v, obteve=%+v", want, got)
			}
		})
	}
}

func TestTokenImmediatelyAfterString(t *testing.T) {
	l := New(`"texto";`)
	want := []Token{
		{Type: STRING_CONST, Literal: "texto", Line: 1},
		{Type: SEMICOLON, Literal: ";", Line: 1},
		{Type: EOF, Literal: "", Line: 1},
	}
	for i, expected := range want {
		if got := l.NextToken(); got != expected {
			t.Fatalf("token[%d] esperado=%+v, obteve=%+v", i, expected, got)
		}
	}
}

func TestUnterminatedStringAtEOF(t *testing.T) {
	for _, content := range []string{"", "texto"} {
		t.Run(content, func(t *testing.T) {
			l := New("\n\"" + content)
			want := Token{Type: ILLEGAL, Literal: content, Line: 2}
			if got := l.NextToken(); got != want {
				t.Fatalf("token esperado=%+v, obteve=%+v", want, got)
			}
			// Chamadas repetidas devem continuar retornando EOF.
			for i := 0; i < 2; i++ {
				want = Token{Type: EOF, Literal: "", Line: 2}
				if got := l.NextToken(); got != want {
					t.Fatalf("apos erro esperado=%+v, obteve=%+v", want, got)
				}
			}
		})
	}
}

func TestStringInterruptedByLineBreak(t *testing.T) {
	tests := []struct {
		name      string
		breakText string
	}{
		{"LF", "\n"},
		{"CR", "\r"},
		{"CRLF", "\r\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New("\n\"texto" + tt.breakText + "x\n\"ok\";")
			want := []Token{
				{Type: ILLEGAL, Literal: "texto", Line: 2},
				{Type: IDENTIFIER, Literal: "x", Line: 3},
				{Type: STRING_CONST, Literal: "ok", Line: 4},
				{Type: SEMICOLON, Literal: ";", Line: 4},
				{Type: EOF, Literal: "", Line: 4},
			}
			for i, expected := range want {
				if got := l.NextToken(); got != expected {
					t.Fatalf("token[%d] esperado=%+v, obteve=%+v", i, expected, got)
				}
			}
		})
	}
}

func TestBackslashDoesNotEscapeQuote(t *testing.T) {
	l := New(`"texto\";`)
	want := []Token{
		{Type: STRING_CONST, Literal: `texto\`, Line: 1},
		{Type: SEMICOLON, Literal: ";", Line: 1},
		{Type: EOF, Literal: "", Line: 1},
	}
	for i, expected := range want {
		if got := l.NextToken(); got != expected {
			t.Fatalf("token[%d] esperado=%+v, obteve=%+v", i, expected, got)
		}
	}
}
