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

	INFIX_OPERATORS_FUNCS = map[string]func(object.Object, object.Object) object.Object{
		token.PLUS:  infixPlus,
		token.MINUS: infixMinus,
		token.TIMES: infixTimes,
		token.SLASH: infixSlash,

		token.EQ:  infixEQ,
		token.NEQ: infixNEQ,
		token.LT:  infixLT,
		token.LEQ: infixLEQ,
		token.GT:  infixGT,
		token.GEQ: infixGEQ,
	}
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

	case *ast.InfixExpression:
		return evalInfixExpression(node)

	}

	return NULL
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

func evalInfixExpression(infixExp *ast.InfixExpression) object.Object {
	r, l := Eval(infixExp.Right), Eval(infixExp.Left)

	infixOp, ok := INFIX_OPERATORS_FUNCS[infixExp.Operator]
	if !ok {
		return NULL
	}

	return infixOp(l, r)
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

func infixPlus(l, r object.Object) object.Object {

	if l == NULL || r == NULL {
		return NULL
	}

	if l.Type() == object.BOOLEAN_OBJ || r.Type() == object.BOOLEAN_OBJ {
		return NULL
	}

	// At this point it is safe to cast
	left, right := l.(*object.Integer), r.(*object.Integer)

	return &object.Integer{
		Value: left.Value + right.Value,
	}
}

func infixMinus(l, r object.Object) object.Object {

	if l == NULL || r == NULL {
		return NULL
	}

	if l.Type() == object.BOOLEAN_OBJ || r.Type() == object.BOOLEAN_OBJ {
		return NULL
	}

	// At this point it is safe to cast
	left, right := l.(*object.Integer), r.(*object.Integer)

	return &object.Integer{
		Value: left.Value - right.Value,
	}
}

func infixTimes(l, r object.Object) object.Object {

	if l == NULL || r == NULL {
		return NULL
	}

	if l.Type() == object.BOOLEAN_OBJ || r.Type() == object.BOOLEAN_OBJ {
		return NULL
	}

	// At this point it is safe to cast
	left, right := l.(*object.Integer), r.(*object.Integer)

	return &object.Integer{
		Value: left.Value * right.Value,
	}
}

func infixSlash(l, r object.Object) object.Object {

	if l == NULL || r == NULL {
		return NULL
	}

	if l.Type() == object.BOOLEAN_OBJ || r.Type() == object.BOOLEAN_OBJ {
		return NULL
	}

	// At this point it is safe to cast
	left, right := l.(*object.Integer), r.(*object.Integer)

	return &object.Integer{
		Value: left.Value / right.Value,
	}
}

func infixEQ(l, r object.Object) object.Object {

	if l == NULL || r == NULL {
		return NULL
	}

	if l.Type() != r.Type() {
		return NULL
	}

	switch r.Type() {
	case object.BOOLEAN_OBJ:
		rBoolean, lBoolean := r.(*object.Boolean).Value, l.(*object.Boolean).Value
		return &object.Boolean{
			Value: lBoolean == rBoolean,
		}

	case object.INTEGER_OBJ:
		rInteger, lInteger := r.(*object.Integer).Value, l.(*object.Integer).Value
		return &object.Boolean{
			Value: lInteger == rInteger,
		}

	default:
		return NULL
	}
}

func infixNEQ(l, r object.Object) object.Object {

	if l == NULL || r == NULL {
		return NULL
	}

	if l.Type() != r.Type() {
		return NULL
	}

	switch r.Type() {
	case object.BOOLEAN_OBJ:
		rBoolean, lBoolean := r.(*object.Boolean).Value, l.(*object.Boolean).Value
		return &object.Boolean{
			Value: lBoolean != rBoolean,
		}

	case object.INTEGER_OBJ:
		rInteger, lInteger := r.(*object.Integer).Value, l.(*object.Integer).Value
		return &object.Boolean{
			Value: lInteger != rInteger,
		}

	default:
		return NULL
	}
}

func infixLEQ(l, r object.Object) object.Object {

	if l == NULL || r == NULL {
		return NULL
	}

	if l.Type() != r.Type() {
		return NULL
	}

	if l.Type() == object.BOOLEAN_OBJ {
		return NULL
	}

	rInteger, lInteger := r.(*object.Integer).Value, l.(*object.Integer).Value
	return &object.Boolean{
		Value: lInteger <= rInteger,
	}

}

func infixLT(l, r object.Object) object.Object {

	if l == NULL || r == NULL {
		return NULL
	}

	if l.Type() != r.Type() {
		return NULL
	}

	if l.Type() == object.BOOLEAN_OBJ {
		return NULL
	}

	rInteger, lInteger := r.(*object.Integer).Value, l.(*object.Integer).Value
	return &object.Boolean{
		Value: lInteger < rInteger,
	}

}

func infixGEQ(l, r object.Object) object.Object {

	if l == NULL || r == NULL {
		return NULL
	}

	if l.Type() != r.Type() {
		return NULL
	}

	if l.Type() == object.BOOLEAN_OBJ {
		return NULL
	}

	rInteger, lInteger := r.(*object.Integer).Value, l.(*object.Integer).Value
	return &object.Boolean{
		Value: lInteger >= rInteger,
	}

}

func infixGT(l, r object.Object) object.Object {

	if l == NULL || r == NULL {
		return NULL
	}

	if l.Type() != r.Type() {
		return NULL
	}

	if l.Type() == object.BOOLEAN_OBJ {
		return NULL
	}

	rInteger, lInteger := r.(*object.Integer).Value, l.(*object.Integer).Value
	return &object.Boolean{
		Value: lInteger > rInteger,
	}

}
