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