package lexer

import "compiler-book/token"

type Lexer interface {
	NextToken() token.Token
}

type lexer struct {
	input        []rune
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current rune under examination
	column       int  // current column in input
	line         int  // current line in input
}

func New(input string) Lexer {
	l := &lexer{input: []rune(input)}
	l.readChar()
	return l
}

func (l *lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Metadata: token.TokenMetadata{Line: l.line, Column: l.column}}

		} else {
			tok = l.newToken(token.ASSIGN, l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Metadata: token.TokenMetadata{Line: l.line, Column: l.column}}
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Metadata: token.TokenMetadata{Line: l.line, Column: l.column}}
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch)
		}
	case ';':
		tok = l.newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = l.newToken(token.LPAREN, l.ch)
	case ')':
		tok = l.newToken(token.RPAREN, l.ch)
	case ',':
		tok = l.newToken(token.COMMA, l.ch)
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INC, Literal: string(ch) + string(l.ch), Metadata: token.TokenMetadata{Line: l.line, Column: l.column}}
		} else {
			tok = l.newToken(token.PLUS, l.ch)
		}
	case '{':
		tok = l.newToken(token.LBRACE, l.ch)
	case '}':
		tok = l.newToken(token.RBRACE, l.ch)
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DEC, Literal: string(ch) + string(l.ch), Metadata: token.TokenMetadata{Line: l.line, Column: l.column}}
		} else {
			tok = l.newToken(token.MINUS, l.ch)
		}
	case '\'':
		tok.Literal, tok.Type = l.readRune()
		tok.Metadata = token.TokenMetadata{Line: l.line, Column: l.column}
	case ':':
		tok = l.newToken(token.COLON, l.ch)
	case '/':
		if l.peekChar() == '/' {
			l.skipComment()
			return l.NextToken()
		}
		tok = l.newToken(token.SLASH, l.ch)
	case '*':
		tok = l.newToken(token.ASTERISK, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Metadata: token.TokenMetadata{Line: l.line, Column: l.column}}
		} else {
			tok = l.newToken(token.BANG, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch), Metadata: token.TokenMetadata{Line: l.line, Column: l.column}}
		} else {
			tok = l.newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch), Metadata: token.TokenMetadata{Line: l.line, Column: l.column}}
		} else {
			tok = l.newToken(token.GT, l.ch)
		}
	case '"':
		tok.Literal, tok.Type = l.readString()
		tok.Metadata = token.TokenMetadata{Line: l.line, Column: l.column}
	case '[':
		tok = l.newToken(token.LBRACKET, l.ch)
	case ']':
		tok = l.newToken(token.RBRACKET, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) { // TODO: isLetter() to support unicode
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		}

		if isDigit(l.ch) { // TODO: isDigit() to support unicode
			tok.Literal, tok.Type = l.readNumber()
			return tok
		}

		tok = l.newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func (l *lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line += 1
			l.column = 0
		}
		l.readChar() // TODO: handle \r\n
	}
}

func isDigit(ch rune) bool { // TODO: isDigit() to support unicode
	return '0' <= ch && ch <= '9'
}

func (l *lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0 // TODO: handle \r\n
	}
	return l.input[l.readPosition]
}

func (l *lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func (l *lexer) readNumber() (string, token.TokenType) {
	position := l.position
	numberType := token.INT

	// handle integers (e.g. 123456)
	for isDigit(l.ch) {
		l.readChar()
	}

	// handle floats (e.g. 1.23456)
	if l.ch == '.' {
		// .123 is not a valid float, so we need to check if the next character is a digit
		numberType = token.FLOAT
		l.readChar()

		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return string(l.input[position:l.position]), numberType
}

func isLetter(ch rune) bool { // TODO: isLetter() to support unicode
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *lexer) newToken(tokenType token.TokenType, ch ...rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch),
		Metadata: token.TokenMetadata{
			Line:   l.line,
			Column: l.column,
		},
	}
}

func (l *lexer) readString() (string, token.TokenType) {
	tokenType := token.ILLEGAL
	position := l.position + 1
	for {
		l.readChar()

		if l.ch == '\\' && l.peekChar() == '"' {
			l.readChar() // skip the escape character
			continue
		}

		if l.ch == '"' {
			tokenType = token.STRING
			break
		}

		if l.ch == 0 {
			break
		}
	}
	return string(l.input[position:l.position]), tokenType
}

func (l *lexer) readRune() (string, token.TokenType) {
	position := l.position + 1

	if position >= len(l.input) { // avoid index out of range
		return "", token.ILLEGAL
	}

	run := l.input[position]

	if l.peekChar() == '\\' {
		l.readChar()
		run = l.input[position+1]
	}

	l.readChar()

	if l.peekChar() != '\'' {
		return string(run), token.ILLEGAL
	}

	l.readChar()

	return string(run), token.RUNE
}

func (l *lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
	l.column += 1
}
