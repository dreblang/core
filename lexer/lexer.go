package lexer

import (
	"strconv"
	"strings"

	"github.com/dreblang/core/token"
)

type Lexer struct {
	input        string
	position     int
	nextPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.Equal, Literal: literal}
		} else {
			tok = newToken(token.Assign, l.ch)
		}
	case '+':
		tok = newToken(token.Plus, l.ch)
	case '-':
		tok = newToken(token.Minus, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NotEqual, Literal: literal}
		} else {
			tok = newToken(token.Bang, l.ch)
		}
	case '*':
		tok = newToken(token.Asterisk, l.ch)
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			tok.Type = token.DoubleSlash
			tok.Literal = l.readComment()
		} else {
			tok = newToken(token.Slash, l.ch)
		}
	case '%':
		tok = newToken(token.Percent, l.ch)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LessOrEqual, Literal: literal}
		} else {
			tok = newToken(token.LessThan, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GreaterOrEqual, Literal: literal}
		} else {
			tok = newToken(token.GreaterThan, l.ch)
		}
	case ',':
		tok = newToken(token.Comma, l.ch)
	case ';':
		tok = newToken(token.Semicolon, l.ch)
	case ':':
		if l.peekChar() == ':' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DoubleColon, Literal: literal}
		} else {
			tok = newToken(token.Colon, l.ch)
		}
	case '.':
		tok = newToken(token.Dot, l.ch)
	case '(':
		tok = newToken(token.LeftParen, l.ch)
	case ')':
		tok = newToken(token.RightParen, l.ch)
	case '{':
		tok = newToken(token.LeftBrace, l.ch)
	case '}':
		tok = newToken(token.RightBrace, l.ch)
	case '[':
		tok = newToken(token.LeftBracket, l.ch)
	case ']':
		tok = newToken(token.RightBracket, l.ch)
	case '"':
		tok.Type = token.String
		tok.Literal = l.readString()
	case '\'':
		tok.Type = token.String
		tok.Literal = l.readString2()
	case '`':
		tok.Type = token.String
		tok.Literal = l.readString3()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) || l.ch == '_' {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifierType(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			literal, isFloat := l.readNumber()
			tok.Literal = literal
			if isFloat {
				tok.Type = token.Float
			} else {
				tok.Type = token.Int
			}
			return tok
		} else {
			tok = newToken(token.Illegal, l.ch)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readChar() {
	l.ch = l.peekChar()
	l.position = l.nextPosition
	l.nextPosition += 1
}

func (l *Lexer) readString() string {
	pos := l.position + 1
	for {
		l.readChar()
		if l.ch == '\\' {
			l.readChar()
			continue
		} else if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	result, _ := strconv.Unquote("\"" + l.input[pos:l.position] + "\"")
	return result
}

func (l *Lexer) readString2() string {
	pos := l.position + 1
	for {
		l.readChar()
		if l.ch == '\\' {
			l.readChar()
			continue
		} else if l.ch == '\'' || l.ch == 0 {
			break
		}
	}
	instr := l.input[pos:l.position]
	instr = strings.ReplaceAll(instr, "\"", "\\\"")
	instr = strings.ReplaceAll(instr, "\\'", "'")
	result, _ := strconv.Unquote("\"" + instr + "\"")
	return result
}

func (l *Lexer) readString3() string {
	pos := l.position + 1
	for {
		l.readChar()
		if l.ch == '\\' {
			l.readChar()
			continue
		} else if l.ch == '`' || l.ch == 0 {
			break
		}
	}
	result, _ := strconv.Unquote("`" + l.input[pos:l.position] + "`")
	return result
}

func (l *Lexer) readComment() string {
	pos := l.position + 1
	for {
		l.readChar()
		if l.ch == '\n' || l.ch == 0 {
			break
		}
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readNumber() (string, bool) {
	seenDot := false
	pos := l.position
	for {
		if isDigit(l.ch) {
			l.readChar()
		} else if l.ch == '.' && !seenDot {
			seenDot = true
			l.readChar()
		} else {
			break
		}
	}
	return l.input[pos:l.position], seenDot
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	if isLetter(l.ch) || l.ch == '_' {
		l.readChar()
	}
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition]
	}
}
