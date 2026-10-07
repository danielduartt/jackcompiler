package lexer

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
		line:  1, // Iniciar na linha 1...
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// Ele olha o próximo caractere sem avançar os ponteiros.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// skipWhitespace avança enquanto encontrar espaços vazios
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++ // Conta as linhas exatamente como no Java
		}
		l.readChar()
	}
}

func newToken(tokenType TokenType, ch byte, line int) Token {
	return Token{Type: tokenType, Literal: string(ch), Line: line}
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	// Símbolos de 1 caractere
	case '=':
		tok = newToken(EQ, l.ch, l.line)
	case '+':
		tok = newToken(PLUS, l.ch, l.line)
	case '-':
		tok = newToken(MINUS, l.ch, l.line)
	case '*':
		tok = newToken(ASTERISK, l.ch, l.line)
	case '<':
		tok = newToken(LT, l.ch, l.line)
	case '>':
		tok = newToken(GT, l.ch, l.line)
	case '&':
		tok = newToken(AND, l.ch, l.line)
	case '|':
		tok = newToken(OR, l.ch, l.line)
	case '~':
		tok = newToken(NOT, l.ch, l.line)
	case '.':
		tok = newToken(DOT, l.ch, l.line)
	case ',':
		tok = newToken(COMMA, l.ch, l.line)
	case ';':
		tok = newToken(SEMICOLON, l.ch, l.line)
	case '{':
		tok = newToken(LBRACE, l.ch, l.line)
	case '}':
		tok = newToken(RBRACE, l.ch, l.line)
	case '(':
		tok = newToken(LPAREN, l.ch, l.line)
	case ')':
		tok = newToken(RPAREN, l.ch, l.line)
	case '[':
		tok = newToken(LBRACKET, l.ch, l.line)
	case ']':
		tok = newToken(RBRACKET, l.ch, l.line)

	// O caso especial da barra '/'
	case '/':
		if l.peekChar() == '/' {
			l.skipLineComment()
			return l.NextToken() // Pula o comentário e tenta o próximo token real
		} else if l.peekChar() == '*' {
			l.skipBlockComment()
			return l.NextToken() // Pula o bloco de comentário e tenta de novo
		} else {
			tok = newToken(SLASH, l.ch, l.line)
		}

	// Tratamento de Strings
	case '"':
		tok.Type = STRING_CONST
		tok.Literal = l.readString()
		tok.Line = l.line
		l.readChar() // Consome as aspas de fechamento da string
		return tok   // Retornamos direto pois a string já avançou os ponteiros

	case 0:
		tok.Literal = ""
		tok.Type = EOF
		tok.Line = l.line

	// O Bloco Default resolve Identificadores, Keywords e Inteiros
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal) // Verifica se é Keyword no map
			tok.Line = l.line
			return tok // Retornamos direto, readIdentifier já deixou o ponteiro no lugar certo
		} else if isDigit(l.ch) {
			tok.Type = INT_CONST
			tok.Literal = l.readNumber()
			tok.Line = l.line
			return tok // Retornamos direto
		} else {
			tok = newToken(ILLEGAL, l.ch, l.line)
		}
	}

	l.readChar() // Avança para o próximo caractere preparatório para símbolos de 1 char
	return tok
}

// ==========================================
// FUNÇÕES AUXILIARES DE LEITURA
// ==========================================

func (l *Lexer) readIdentifier() string {
	position := l.position
	// Um identificador em Jack pode começar com letra/underscore, mas depois pode conter números
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1 // Pula a aspa dupla de abertura
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	// Deixamos o \n para ser consumido e contado no skipWhitespace da próxima rodada
}

func (l *Lexer) skipBlockComment() {
	l.readChar() // Consome o '/'
	l.readChar() // Consome o '*'
	for {
		if l.ch == 0 {
			break // Fim do arquivo inesperado
		}
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // Consome '*'
			l.readChar() // Consome '/'
			l.readChar() // Avança para o próximo caractere após o comentário
			break
		}
		if l.ch == '\n' {
			l.line++ // Conta quebras de linha dentro de comentários de bloco
		}
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
