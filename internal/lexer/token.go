package lexer

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Literais e Identificadores
	IDENTIFIER   = "IDENTIFIER"
	INT_CONST    = "INT_CONST"
	STRING_CONST = "STRING_CONST"

	// Símbolos
	PLUS      = "+"
	MINUS     = "-"
	ASTERISK  = "*"
	SLASH     = "/"
	DOT       = "."
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
	AND       = "&"
	OR        = "|"
	NOT       = "~"
	LT        = "<"
	GT        = ">"
	EQ        = "="
)

// 2. Tabela de Palavras-chave
var keywords = map[string]TokenType{
	"class":       "CLASS",
	"constructor": "CONSTRUCTOR",
	"function":    "FUNCTION",
	"method":      "METHOD",
	"field":       "FIELD",
	"static":      "STATIC",
	"var":         "VAR",
	"int":         "INT",
	"char":        "CHAR",
	"boolean":     "BOOLEAN",
	"void":        "VOID",
	"true":        "TRUE",
	"false":       "FALSE",
	"null":        "NULL",
	"this":        "THIS",
	"let":         "LET",
	"do":          "DO",
	"if":          "IF",
	"else":        "ELSE",
	"while":       "WHILE",
	"return":      "RETURN",
}

// LookupIdent verifica se um identificador é uma palavra-chave.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}
