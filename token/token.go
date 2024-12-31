package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

// Define different token types in the language

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPARENTHESIS = "("
	RPARENTHESIS = ")"
	LBRACE       = "{"
	RBRACE       = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

// Helper function that takes as parameter a token type and its corresponding char in the case of EOF, operators, delimiters, parenthesis and braces, and returns the corresponding Token instance.
func NewToken(tokenType TokenType, ch byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

// Helper function, takes as parameter a byte, representing a one char token, and returns its corresponding token type.
func CharToToken(ch byte) TokenType {

	var tt TokenType

	switch ch {
	case '=':
		tt = ASSIGN
	case '+':
		tt = PLUS
	case ',':
		tt = COMMA
	case ';':
		tt = SEMICOLON
	case '(':
		tt = LPARENTHESIS
	case ')':
		tt = RPARENTHESIS
	case '{':
		tt = LBRACE
	case '}':
		tt = RBRACE
	case 0:
		tt = EOF
	default:
		// TODO: Not sure about this step. Will certainly change it afterwards once the lexer takes into account identifiers and keywords ...
		tt = ILLEGAL
	}

	return tt
}
