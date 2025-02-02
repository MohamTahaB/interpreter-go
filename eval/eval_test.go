package eval

import (
	"testing"

	"github.com/MohamTahaB/interpreter-go/lexer"
	"github.com/MohamTahaB/interpreter-go/object"
	"github.com/MohamTahaB/interpreter-go/parser"
)

func TestIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"69", 69},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not of type Integer. Got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. Expect=%d, got=%d", expected, result.Value)
		return false
	}

	return true

}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()

	return Eval(program)
}
