package scanner

type TokenType int

// Estrutura de um Token <tipo, valor>
type Token struct {
	Type  TokenType
	Value string
	Line  int // linha para itentificar erros futuramente
}

const (
	// Em Go, toda variável nasce com o valor zero do seu tipo, então Se o 0 fosse LPAREN,
	// um bug de inicialização pareceria um parêntese legítimo. Reservando o 0 para ILLEGAL,
	// o erro fica evidente.
	ILLEGAL TokenType = iota
	// Símbolos (19)
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
	LT
	GT
	EQ

	// Literais e identificadores
	NUMBER
	IDENT
	STRING

	// Palavras reservadas (21)
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

	// Fim da entrada. O XML de teste NÃO imprime este token.
	EOF
)

// Aqui a gente faz como se fosse um dicionário em python, pq no nosso "enum" esses valores são associados a um valor int
// então se eu quiser saber o nome de um token, posso usar este array como um dicionário para verificar o nome do tipo e não seu inteiro
// parece um dicionário mais não é. Aqui é um array indexado pelo nome das constantes do TokenType.... (meio doido)
var tokenNames = [...]string{
	ILLEGAL:     "ILLEGAL",
	LPAREN:      "LPAREN",
	RPAREN:      "RPAREN",
	LBRACE:      "LBRACE",
	RBRACE:      "RBRACE",
	LBRACKET:    "LBRACKET",
	RBRACKET:    "RBRACKET",
	COMMA:       "COMMA",
	SEMICOLON:   "SEMICOLON",
	DOT:         "DOT",
	PLUS:        "PLUS",
	MINUS:       "MINUS",
	ASTERISK:    "ASTERISK",
	SLASH:       "SLASH",
	AND:         "AND",
	OR:          "OR",
	NOT:         "NOT",
	LT:          "LT",
	GT:          "GT",
	EQ:          "EQ",
	NUMBER:      "NUMBER",
	IDENT:       "IDENT",
	STRING:      "STRING",
	CLASS:       "CLASS",
	CONSTRUCTOR: "CONSTRUCTOR",
	FUNCTION:    "FUNCTION",
	METHOD:      "METHOD",
	FIELD:       "FIELD",
	STATIC:      "STATIC",
	VAR:         "VAR",
	INT:         "INT",
	CHAR:        "CHAR",
	BOOLEAN:     "BOOLEAN",
	VOID:        "VOID",
	TRUE:        "TRUE",
	FALSE:       "FALSE",
	NULL:        "NULL",
	THIS:        "THIS",
	LET:         "LET",
	DO:          "DO",
	IF:          "IF",
	ELSE:        "ELSE",
	WHILE:       "WHILE",
	RETURN:      "RETURN",
	EOF:         "EOF",
}

// FUNÇÃO PARA DEPURAR: devolve o nome do token como string <ILLEGAL, LPAREN, RPAREN, ...>
func (t TokenType) String() string {
	// Sem esta checagem, um valor fora do intervalo causaria panic.
	if t < 0 || int(t) >= len(tokenNames) || tokenNames[t] == "" {
		return "UNKNOWN"
	}
	return tokenNames[t]
}

// keywords mapeia cada palavra reservada de Jack ao seu tipo de token.
var keywords = map[string]TokenType{
	"class":       CLASS,
	"constructor": CONSTRUCTOR,
	"function":    FUNCTION,
	"method":      METHOD,
	"field":       FIELD,
	"static":      STATIC,
	"var":         VAR,
	"int":         INT,
	"char":        CHAR,
	"boolean":     BOOLEAN,
	"void":        VOID,
	"true":        TRUE,
	"false":       FALSE,
	"null":        NULL,
	"this":        THIS,
	"let":         LET,
	"do":          DO,
	"if":          IF,
	"else":        ELSE,
	"while":       WHILE,
	"return":      RETURN,
}

// symbols mapeia cada símbolo de Jack (um único caractere) ao seu tipo.
// A chave é rune porque o scanner vai percorrer um []rune.
var symbols = map[rune]TokenType{
	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,
	'[': LBRACKET,
	']': RBRACKET,
	',': COMMA,
	';': SEMICOLON,
	'.': DOT,
	'+': PLUS,
	'-': MINUS,
	'*': ASTERISK,
	'/': SLASH,
	'&': AND,
	'|': OR,
	'~': NOT,
	'<': LT,
	'>': GT,
	'=': EQ,
}

/* FORMA ANTIGA DE CATEGORY
// Category devolve o nome da tag XML do nand2tetris para este tipo de token
// Devolve "" para ILLEGAL e EOF, que não aparecem no XML.
func (t TokenType) Category() string {
	switch {
	case t >= LPAREN && t <= EQ:
		return "symbol"
	case t >= CLASS && t <= RETURN:
		return "keyword"
	case t == NUMBER:
		return "integerConstant"
	case t == STRING:
		return "stringConstant"
	case t == IDENT:
		return "identifier"
	}
	return ""
}
*/

// FORMA NOVA SEM PROBLEMAS (PELO MENOS EU ACHO)
var categories = map[TokenType]string{}

func init() {
	for _, t := range keywords {
		categories[t] = "keyword"
	}
	for _, t := range symbols {
		categories[t] = "symbol"
	}
	categories[NUMBER] = "integerConstant"
	categories[STRING] = "stringConstant"
	categories[IDENT] = "identifier"
}
func (t TokenType) Category() string {
	return categories[t] // tipo ausente devolve "" (serve para ILLEGAL e EOF)
}

// a Cadeia ou Palavra é um identificador ou é reservada? Essa função define isso
func lookupIdent(word string) TokenType {
	// comma-ok: "ok" distingue "chave ausente" de "chave presente com valor zero".
	if tipo, ok := keywords[word]; ok {
		return tipo
	}
	return IDENT
}
