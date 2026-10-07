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
