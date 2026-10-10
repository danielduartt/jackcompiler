package scanner

import (
	"fmt"
	"strconv"
)

const maxInt = 32767

type Scanner struct {
	input []rune // []rune: cada posição é um caractere inteiro (Unicode-safe)
	pos   int    // índice do próximo caractere a ser lido
	line  int    // linha atual (começa em 1)
}

// Por conversão a função contrutora em Go se chamada new
func New(src string) *Scanner {
	return &Scanner{input: []rune(src), line: 1}
}

// Navegação na cadeia de caracteres
func (s *Scanner) atEnd() bool {
	return s.pos >= len(s.input)
}

// peek vai olhar o caractere atual sem avançar a posição, ou seja sem consumir. Vai devolver 0 se for o fim da entrada.
func (s *Scanner) peek() rune {
	if s.atEnd() {
		return 0
	}
	return s.input[s.pos]
}

// peekNext vai olhar o próximo caractere sem avançar a posição. Vai seguir a mesma lógica do peek
func (s *Scanner) peekNext() rune {
	if s.pos+1 >= len(s.input) {
		return 0
	}
	return s.input[s.pos+1]
}

// advance vai avançar a posição consumindo o caractere atual e retornando ele
func (s *Scanner) advance() rune {
	if s.atEnd() {
		return 0
	}
	r := s.input[s.pos]
	s.pos++
	// aqui utilizamos a condicional para atualizar a linha atual quando encontramos um caractere de nova linha
	// conta quebras de linhas
	if r == '\n' {
		s.line++
	}
	return r
}

// Classificação de caracteres

func isDigit(r rune) bool  { return r >= '0' && r <= '9' }
func isLetter(r rune) bool { return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' }

// Ruídos, Espaços e Comentários são limpos pelo skipIgnored que vai pular os espaços, tabs e etc...
// aqui tem como devolver um erro, se um comentário de bloco por exemplo ficar sem fechar
func (s *Scanner) skipIgnored() error {
	for !s.atEnd() {
		r := s.peek()
		switch {
		// caso queiramos ignorar espaços em branco
		case r == ' ' || r == '\t' || r == '\r' || r == '\n':
			s.advance()
		// caso queiramos ignorar comentários de linha um / seguido de /
		case r == '/' && s.peekNext() == '/':
			for !s.atEnd() && s.peek() != '\n' {
				s.advance()
			}
		case r == '/' && s.peekNext() == '*': // comentário de bloco (inclui /** */)
			startLine := s.line
			s.advance() // '/'
			s.advance() // '*'
			closed := false
			for !s.atEnd() {
				if s.peek() == '*' && s.peekNext() == '/' {
					s.advance()
					s.advance()
					closed = true
					break // em Go, break sai do for mais interno (não do switch)
				}
				s.advance()
			}
			if !closed {
				return fmt.Errorf("linha %d: comentário de bloco não fechado", startLine)
			}
		default:
			return nil
		}
	}
	return nil
}

// Leitura dos tokens

func (s *Scanner) readNumber(line int) (Token, error) {
	start := s.pos
	for !s.atEnd() && isDigit(s.peek()) {
		s.advance()
	}
	text := string(s.input[start:s.pos]) // fatia [início:fim), como em Python
	n, err := strconv.Atoi(text)
	if err != nil || n > maxInt {
		return Token{Type: ILLEGAL, Value: text, Line: line},
			fmt.Errorf("linha %d: inteiro fora do intervalo 0..%d: %s", line, maxInt, text)
	}
	return Token{Type: NUMBER, Value: text, Line: line}, nil
}

func (s *Scanner) readString(line int) (Token, error) {
	s.advance() // aspa de abertura (não faz parte do valor)
	start := s.pos
	for !s.atEnd() && s.peek() != '"' && s.peek() != '\n' {
		s.advance()
	}
	if s.atEnd() || s.peek() == '\n' {
		return Token{Type: ILLEGAL, Line: line},
			fmt.Errorf("linha %d: string não fechada", line)
	}
	value := string(s.input[start:s.pos])
	s.advance() // aspa de fechamento
	return Token{Type: STRING, Value: value, Line: line}, nil
}

func (s *Scanner) readWord(line int) Token {
	start := s.pos
	for !s.atEnd() && (isLetter(s.peek()) || isDigit(s.peek())) {
		s.advance()
	}
	word := string(s.input[start:s.pos])
	return Token{Type: lookupIdent(word), Value: word, Line: line}
}

func (s *Scanner) NextToken() (Token, error) {
	if err := s.skipIgnored(); err != nil {
		return Token{Type: ILLEGAL, Line: s.line}, err
	}
	if s.atEnd() {
		return Token{Type: EOF, Line: s.line}, nil
	}

	line := s.line // linha em que o token COMEÇA
	r := s.peek()

	switch {
	case r == '"':
		return s.readString(line)
	case isDigit(r):
		return s.readNumber(line)
	case isLetter(r):
		return s.readWord(line), nil
	}

	if tipo, ok := symbols[r]; ok {
		s.advance()
		return Token{Type: tipo, Value: string(r), Line: line}, nil
	}

	s.advance()
	return Token{Type: ILLEGAL, Value: string(r), Line: line},
		fmt.Errorf("linha %d: caractere inesperado %q", line, r)
}

func (s *Scanner) Tokens() ([]Token, error) {
	var tokens []Token
	for {
		tok, err := s.NextToken()
		if err != nil {
			return tokens, err
		}
		if tok.Type == EOF {
			return tokens, nil
		}
		tokens = append(tokens, tok)
	}
}
