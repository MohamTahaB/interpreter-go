package eval

import (
	"github.com/MohamTahaB/interpreter-go/ast"
	"github.com/MohamTahaB/interpreter-go/object"
	"github.com/MohamTahaB/interpreter-go/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:

		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	}

	return nil
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)
	}

	return result
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {

	case token.NEG:
		return evalNegationPrefixExpression(right)
	case token.MINUS:
		return evalNegativePrefixExpression(right)

	default:
		return NULL
	}
}

func evalNegationPrefixExpression(right object.Object) object.Object {
	var nativeRightEval bool

	if right == NULL {
		return NULL
	}

	// Native prefix eval to be determined
	switch prefixEval := right.(type) {
	case *object.Boolean:
		nativeRightEval = booleanObjectToNativeBool(prefixEval)
	case *object.Integer:
		nativeRightEval = prefixEval.Value != 0
	}
	return nativeBoolToBooleanObject(!nativeRightEval)
}

func evalNegativePrefixExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	val := right.(*object.Integer).Value
	return &object.Integer{Value: -val}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

func booleanObjectToNativeBool(input *object.Boolean) bool {
	return input == TRUE
}
