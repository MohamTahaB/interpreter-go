package lexer

import (
	"github.com/MohamTahaB/interpreter-go/token"
)

type Lexer struct {
	input        string // The code being lexed.
	position     int    // Current pos in input, points to curr char
	readPosition int    // Current reading pos, points to next char
	ch           byte   // Current char
}

// Lexer attributes are more or less self explanatory. The reason why we have two pointers: position and readPosition, is that we will need to peek further into the input to see what comes up next

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Helper function to update the position of the considered char in the Lexer instance.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

// Helper function to peek into the readPosition char without reading.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// Returns the next token the Lexer instance points to.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch {
	case token.LegalOneCharLiteral(l.ch):

		var tokType token.TokenType
		var tokenBytes []byte = []byte{l.ch}
		twoCharsToken := true

		switch {
		case l.ch == '!' && l.peekChar() == '=':
			tokType = token.NEQ
		case l.ch == '=' && l.peekChar() == '=':
			tokType = token.EQ
		case l.ch == '<' && l.peekChar() == '=':
			tokType = token.LEQ
		case l.ch == '>' && l.peekChar() == '=':
			tokType = token.GEQ
		case l.ch == '+' && l.peekChar() == '=':
			tokType = token.PLUSEQ
		case l.ch == '-' && l.peekChar() == '=':
			tokType = token.MINUSEQ
		case l.ch == '*' && l.peekChar() == '=':
			tokType = token.TIMESEQ
		case l.ch == '/' && l.peekChar() == '=':
			tokType = token.SLASHEQ
		default:
			tokType = token.CharToToken(l.ch)
			twoCharsToken = false
		}

		// Construct the token.
		if twoCharsToken {
			l.readChar()
			tokenBytes = append(tokenBytes, l.ch)
		}
		tok = token.NewToken(tokType, tokenBytes)
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		}
		if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		}

		if l.ch == '"' {
			tok.Type = token.STRING
			tok.Literal = l.readString()
		} else {
			tok = token.NewToken(token.ILLEGAL, []byte{l.ch})
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
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
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
