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

	// Operators (one char)
	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"
	SLASH  = "/"
	TIMES  = "*"
	LT     = "<"
	GT     = ">"
	NEG    = "!"

	// Operators (two chars)
	EQ      = "=="
	NEQ     = "!="
	LEQ     = "<="
	GEQ     = ">="
	PLUSEQ  = "+="
	MINUSEQ = "-="
	SLASHEQ = "/="
	TIMESEQ = "*="

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
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	// Booleans
	TRUE  = "TRUE"
	FALSE = "FALSE"
)

var ONE_CHAR_TOKEN_LITTERALS map[byte]bool = map[byte]bool{
	'=': true,
	'+': true,
	'-': true,
	'/': true,
	'*': true,
	'<': true,
	'>': true,
	'!': true,
	';': true,
	',': true,
	'(': true,
	')': true,
	'{': true,
	'}': true,
	0:   true,
}

// Helper function that takes as parameter a token type and its corresponding char slice, and returns the corresponding Token instance.
func NewToken(tokenType TokenType, ch []byte) Token {
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
	case '-':
		tt = MINUS
	case '*':
		tt = TIMES
	case '<':
		tt = LT
	case '>':
		tt = GT
	case '!':
		tt = NEG
	case '/':
		tt = SLASH
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

func LegalOneCharLiteral(ch byte) bool {
	return ONE_CHAR_TOKEN_LITTERALS[ch]
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
