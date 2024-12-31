package lexer

import "github.com/MohamTahaB/interpreter-go/token"

type Lexer struct {
	input        string // The code being lexed.
	position     int    // Current pos in input, points to curr char
	readPosition int    // Current reading pos, points to next char
	ch           byte   // Current char
}

// Lexer attributes are more or less self explanatory. The reason why we have two pointers: position and readPosition, is that we will need to peek further into the input to see what comes up next

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.ReadChar()
	return l
}

// Helper function to update the position of the considered char in the Lexer instance.
func (l *Lexer) ReadChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

// Returns the next token the Lexer instance points to.
func (l *Lexer) NextToken() token.Token {
	tok := token.NewToken(token.CharToToken(l.ch), l.ch)
	l.ReadChar()
	return tok
}
