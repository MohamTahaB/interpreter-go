package lexer

import (
	"testing"

	"github.com/MohamTahaB/interpreter-go/token"
)

type testStruct struct {
	expectedType    token.TokenType
	expectedLiteral string
}

// Test one char tokens.
func TestNextToken_oneChar_OK(t *testing.T) {
	input := `=+(){},;`

	tests := []testStruct{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPARENTHESIS, "("},
		{token.RPARENTHESIS, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, "\x00"},
	}

	l := New(input)

	// Go through tests
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// Test multiple char tokens.
func TestNextToken_OK(t *testing.T) {

	input := `let five = 5;
	let ten = 10;

	let add = fn(x, y) {
	x + y;
	};

	let result = add(five, ten);
	!-/*5;
	5 < 10 > 5;
	`

	l := New(input)

	tests := []testStruct{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPARENTHESIS, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPARENTHESIS, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPARENTHESIS, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPARENTHESIS, ")"},
		{token.SEMICOLON, ";"},
		{token.NEG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.TIMES, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.EOF, "\x00"},
	}

	for i, tt := range tests {

		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}

}
