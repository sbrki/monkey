package lexer

import "github.com/sbrki/monkey/pkg/token"

type Lexer struct {
	input    string
	currPos  int
	currChar byte
	nextPos  int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar advances the lexer by one character.
func (l *Lexer) readChar() {
	if l.nextPos >= len(l.input) {
		l.currChar = 0 // EOF
	} else {
		l.currChar = l.input[l.nextPos]
	}
	l.currPos = l.nextPos
	l.nextPos += 1
}

func newToken(t token.TokenType, ch byte) token.Token {
	return token.Token{Type: t, Literal: string(ch)}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.currChar {
	case '=':
		tok = newToken(token.ASSIGN, l.currChar)
	case ';':
		tok = newToken(token.SEMICOLON, l.currChar)
	case '(':
		tok = newToken(token.LPAREN, l.currChar)
	case ')':
		tok = newToken(token.RPAREN, l.currChar)
	case ',':
		tok = newToken(token.COMMA, l.currChar)
	case '+':
		tok = newToken(token.PLUS, l.currChar)
	case '{':
		tok = newToken(token.LBRACE, l.currChar)
	case '}':
		tok = newToken(token.RBRACE, l.currChar)
	case 0:
		tok = newToken(token.EOF, 0)
	}

	l.readChar()
	return tok
}
