package ast

import (
	"testing"

	"github.com/MohamTahaB/interpreter-go/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "var"},
					Value: "var",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "another_var"},
					Value: "another_var",
				},
			},
		},
	}

	if program.String() != "let var = another_var;" {
		t.Errorf("program.String() wrong. Got=%q", program.String())
	}
}
