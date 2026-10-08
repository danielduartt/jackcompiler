package scanner

type TokenType int 

const (
	KEYWORD TokenType = iota 
    LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET
    COMMA
	SEMICOLON
	DOT
	PLUS
	MINUS
	ASTERISK
	SLASH
    AND
	OR
	NOT
	LTGT
	EQ
    NUMBER
	IDENT
	STRING
    CLASS
	CONSTRUCTOR
	FUNCTION
	METHOD
	FIELD
	STATIC
	VAR
    INT
	CHAR
	BOOLEAN
	VOID
	TRUE
	FALSE
	NULL
	THIS
    LET
	DO
	IF
	ELSE
	WHILE
	RETURN
	EOF
)

type Token struct {
	Type TokenType
	Value string 
	Line int // linha para itentificar erros futuramente
}

//como se fosse um método (Em Go é meio esquisito isso, mas funciona até que legal)
func (t TokenType) String() string {
	switch t {
	case KEYWORD:
		return "KEYWORD"
	case SYMBOL:
		return "SYMBOL"
	case INT_CONST:
		return "INT_CONST"
	case STRING_CONST:
		return "STRING_CONST"
	case IDENTIFIER:
		return "IDENTIFIER"
	default:
		return "UNKNOWN"
	}
}

// map com as palavras reservadas do Jack 
var keywords = map[string]bool{
	"class" : true, 
	"constructor" : true,
	"function" : true,
	"method" : true,
	"field" : true,
	"static" : true,
	"var" : true, 
	"int" : true,
	"char" : true,
	"boolean" : true,
	"void" : true,
	"true" : true,
	"false" : true,
	"null" : true,
	"this" : true,
	"let" : true,
	"do" : true,
	"idf": true,
	"else": true,
	"while": true,
	"return": true,
}