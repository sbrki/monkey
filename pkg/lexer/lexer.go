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

func (l *Lexer) peekChar() byte {
	if l.nextPos >= len(l.input) {
		return 0 // EOF
	}
	return l.input[l.nextPos]
}

func newToken(t token.TokenType, ch byte) token.Token {
	return token.Token{Type: t, Literal: string(ch)}
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' ||
		'A' <= char && char <= 'Z' ||
		char == '_'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func isWhitespace(char byte) bool {
	return char == ' ' ||
		char == '\t' ||
		char == '\n' ||
		char == '\r'
}

func (l *Lexer) readIdentifier() string {
	startPos := l.currPos
	for isLetter(l.currChar) {
		l.readChar()
	}
	return l.input[startPos:l.currPos]
}

func (l *Lexer) readNumber() string {
	startPos := l.currPos
	for isDigit(l.currChar) {
		l.readChar()
	}
	return l.input[startPos:l.currPos]
}

func (l *Lexer) readString() string {
	//TODO(sbrki): support for basic escape sequences:
	// \t, \n, \', \", \\ 
	l.readChar()
	startPos := l.currPos
	for l.currChar != '"' && l.currChar != 0 {
		l.readChar()
	}
	return l.input[startPos:l.currPos]
}

func (l *Lexer) consumeWhitespace() {
	for isWhitespace(l.currChar) {
		l.readChar()
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.consumeWhitespace()

	switch l.currChar {
	case '=':
		if l.peekChar() == '=' {
			ch := l.currChar
			l.readChar()
			literal := string(ch) + string(l.currChar)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.currChar)
		}
	case '+':
		tok = newToken(token.PLUS, l.currChar)
	case '-':
		tok = newToken(token.MINUS, l.currChar)
	case '!':
		if l.peekChar() == '=' {
			ch := l.currChar
			l.readChar()
			literal := string(ch) + string(l.currChar)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.currChar)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.currChar)
	case '/':
		tok = newToken(token.SLASH, l.currChar)
	case '<':
		tok = newToken(token.LT, l.currChar)
	case '>':
		tok = newToken(token.GT, l.currChar)
	case ',':
		tok = newToken(token.COMMA, l.currChar)
	case ';':
		tok = newToken(token.SEMICOLON, l.currChar)
	case '(':
		tok = newToken(token.LPAREN, l.currChar)
	case ')':
		tok = newToken(token.RPAREN, l.currChar)
	case '{':
		tok = newToken(token.LBRACE, l.currChar)
	case '}':
		tok = newToken(token.RBRACE, l.currChar)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok = newToken(token.EOF, 0)

	default:
		if isLetter(l.currChar) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok // don't advance the lexer! readIdentifier stops at first non-letter char
		} else if isDigit(l.currChar) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok // don't advance the lexer! readNumber stops at first non-num char
		} else {
			tok = newToken(token.ILLEGAL, l.currChar)
		}
	}

	l.readChar()
	return tok
}
