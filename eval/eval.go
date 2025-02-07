package eval

import (
	"fmt"

	"github.com/MohamTahaB/interpreter-go/ast"
	"github.com/MohamTahaB/interpreter-go/object"
	"github.com/MohamTahaB/interpreter-go/token"
)

const (
	UNKNOWN_OP_PREFIX_MSG   = "unknown operator: %s%s"
	UNKNOWN_OP_INFIX_MSG    = "unknown operator: %s %s %s"
	TYPE_MISMATCH_INFIX_MSG = "type mismatch: %s %s %s"
	DIVISION_BY_ZERO        = "division by 0"
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
		return evalProgram(node.Statements)

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
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:

		left := Eval(node.Left)
		if isError(left) {
			return left
		}

		right := Eval(node.Left)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node)

	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements)

	case *ast.IfExpression:
		return evalConditionalExpression(node)

	case *ast.ReturnStatement:
		return &object.ReturnValue{
			Value: Eval(node.ReturnValue),
		}

	}

	return NULL
}

func evalProgram(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)

		if returnVal, ok := result.(*object.ReturnValue); ok {
			return returnVal.Value
		}
		if errorObj, ok := result.(*object.Error); ok {
			return errorObj
		}
	}

	return result
}

func evalBlockStatement(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)

		if returnVal, ok := result.(*object.ReturnValue); ok && result != nil {
			return returnVal
		}
		if errorObj, ok := result.(*object.Error); ok && result != nil {
			return errorObj
		}
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
		return newError(UNKNOWN_OP_PREFIX_MSG, op, right.Type())
	}
}

func evalInfixExpression(infixExp *ast.InfixExpression) object.Object {
	r, l := Eval(infixExp.Right), Eval(infixExp.Left)

	infixOp, ok := INFIX_OPERATORS_FUNCS[infixExp.Operator]
	if !ok {
		// What is not recognized is the operator
		return newError(UNKNOWN_OP_INFIX_MSG, "", infixExp.Operator, "")
	}

	return infixOp(l, r)
}

func evalConditionalExpression(conditionalExp *ast.IfExpression) object.Object {
	conditionEval := Eval(conditionalExp.Condition)

	if isError(conditionEval) {
		return conditionEval
	}

	// Truthy == not null
	if conditionEval.Truthy() {
		return Eval(conditionalExp.Consequence)
	}
	if conditionalExp.Alternative != nil {
		return Eval(conditionalExp.Alternative)
	}

	return NULL
}

func evalNegationPrefixExpression(right object.Object) object.Object {
	var nativeRightEval bool

	if right == NULL {
		return newError(UNKNOWN_OP_PREFIX_MSG, token.NEG, right.Type())
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
		return newError(UNKNOWN_OP_PREFIX_MSG, token.MINUS, right.Type())
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

	if l.Type() != object.INTEGER_OBJ || r.Type() != object.INTEGER_OBJ {

		// Whether a mismatch or unknow op
		if l.Type() == r.Type() {
			return newError(UNKNOWN_OP_INFIX_MSG, l.Type(), token.PLUS, r.Type())
		}

		return newError(TYPE_MISMATCH_INFIX_MSG, l.Type(), token.PLUS, r.Type())
	}

	// At this point it is safe to cast
	left, right := l.(*object.Integer), r.(*object.Integer)

	return &object.Integer{
		Value: left.Value + right.Value,
	}
}

func infixMinus(l, r object.Object) object.Object {

	if l.Type() != object.INTEGER_OBJ || r.Type() != object.INTEGER_OBJ {

		// Whether a mismatch or unknow op
		if l.Type() == r.Type() {
			return newError(UNKNOWN_OP_INFIX_MSG, l.Type(), token.MINUS, r.Type())
		}

		return newError(TYPE_MISMATCH_INFIX_MSG, l.Type(), token.MINUS, r.Type())
	}

	// At this point it is safe to cast
	left, right := l.(*object.Integer), r.(*object.Integer)

	return &object.Integer{
		Value: left.Value - right.Value,
	}
}

func infixTimes(l, r object.Object) object.Object {

	if l.Type() != object.INTEGER_OBJ || r.Type() != object.INTEGER_OBJ {

		// Whether a mismatch or unknow op
		if l.Type() == r.Type() {
			return newError(UNKNOWN_OP_INFIX_MSG, l.Type(), token.TIMES, r.Type())
		}

		return newError(TYPE_MISMATCH_INFIX_MSG, l.Type(), token.TIMES, r.Type())
	}

	// At this point it is safe to cast
	left, right := l.(*object.Integer), r.(*object.Integer)

	return &object.Integer{
		Value: left.Value * right.Value,
	}
}

func infixSlash(l, r object.Object) object.Object {

	if l.Type() != object.INTEGER_OBJ || r.Type() != object.INTEGER_OBJ {

		// Whether a mismatch or unknow op
		if l.Type() == r.Type() {
			return newError(UNKNOWN_OP_INFIX_MSG, l.Type(), token.SLASH, r.Type())
		}

		return newError(TYPE_MISMATCH_INFIX_MSG, l.Type(), token.SLASH, r.Type())
	}

	// At this point it is safe to cast
	left, right := l.(*object.Integer), r.(*object.Integer)

	if right.Value == 0 {
		return newError(DIVISION_BY_ZERO)
	}

	return &object.Integer{
		Value: left.Value / right.Value,
	}
}

func infixEQ(l, r object.Object) object.Object {

	if l == NULL || r == NULL {
		return &object.Boolean{
			Value: l == r,
		}
	}

	if l.Type() != r.Type() {
		return newError(UNKNOWN_OP_INFIX_MSG, l.Type(), token.EQ, r.Type())
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
		return newError(UNKNOWN_OP_INFIX_MSG, l.Type(), token.EQ, r.Type())
	}
}

func infixNEQ(l, r object.Object) object.Object {

	if l == NULL || r == NULL {
		return &object.Boolean{
			Value: l != r,
		}
	}

	if l.Type() != r.Type() {
		return newError(UNKNOWN_OP_INFIX_MSG, l.Type(), token.NEQ, r.Type())
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
		return newError(UNKNOWN_OP_INFIX_MSG, l.Type(), token.NEQ, r.Type())
	}
}

func infixLEQ(l, r object.Object) object.Object {

	if l.Type() != object.INTEGER_OBJ || r.Type() != object.INTEGER_OBJ {
		return newError(UNKNOWN_OP_INFIX_MSG, l.Type(), token.LEQ, r.Type())
	}

	rInteger, lInteger := r.(*object.Integer).Value, l.(*object.Integer).Value
	return &object.Boolean{
		Value: lInteger <= rInteger,
	}

}

func infixLT(l, r object.Object) object.Object {

	if l.Type() != object.INTEGER_OBJ || r.Type() != object.INTEGER_OBJ {
		return newError(UNKNOWN_OP_INFIX_MSG, l.Type(), token.LT, r.Type())
	}

	rInteger, lInteger := r.(*object.Integer).Value, l.(*object.Integer).Value
	return &object.Boolean{
		Value: lInteger < rInteger,
	}

}

func infixGEQ(l, r object.Object) object.Object {

	if l.Type() != object.INTEGER_OBJ || r.Type() != object.INTEGER_OBJ {
		return newError(UNKNOWN_OP_INFIX_MSG, l.Type(), token.GEQ, r.Type())
	}

	rInteger, lInteger := r.(*object.Integer).Value, l.(*object.Integer).Value
	return &object.Boolean{
		Value: lInteger >= rInteger,
	}

}

func infixGT(l, r object.Object) object.Object {

	if l.Type() != object.INTEGER_OBJ || r.Type() != object.INTEGER_OBJ {
		return newError(UNKNOWN_OP_INFIX_MSG, l.Type(), token.GT, r.Type())
	}

	rInteger, lInteger := r.(*object.Integer).Value, l.(*object.Integer).Value
	return &object.Boolean{
		Value: lInteger > rInteger,
	}

}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj == nil {
		return false
	}

	_, ok := obj.(*object.Error)
	return ok
}
