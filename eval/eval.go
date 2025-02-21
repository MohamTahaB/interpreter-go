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
	IDENT_NOT_FOUND         = "identifier not found: %s"
	NOT_A_FUNC              = "not a function: %s"
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

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:

		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(left, right, node.Operator)

	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.IfExpression:
		return evalConditionalExpression(node, env)

	case *ast.ReturnStatement:
		return &object.ReturnValue{
			Value: Eval(node.ReturnValue, env),
		}

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{
			Parameters: params,
			Env:        env,
			Body:       body,
		}

	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(fn, args)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	}

	return NULL
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)

		if returnVal, ok := result.(*object.ReturnValue); ok {
			return returnVal.Value
		}
		if errorObj, ok := result.(*object.Error); ok {
			return errorObj
		}
	}

	return result
}

func evalBlockStatement(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)

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

func evalInfixExpression(l, r object.Object, operator string) object.Object {

	infixOp, ok := INFIX_OPERATORS_FUNCS[operator]
	if !ok {
		// What is not recognized is the operator
		return newError(UNKNOWN_OP_INFIX_MSG, "", operator, "")
	}

	return infixOp(l, r)
}

func evalConditionalExpression(conditionalExp *ast.IfExpression, env *object.Environment) object.Object {
	conditionEval := Eval(conditionalExp.Condition, env)

	if isError(conditionEval) {
		return conditionEval
	}

	// Truthy == not null
	if conditionEval.Truthy() {
		return Eval(conditionalExp.Consequence, env)
	}
	if conditionalExp.Alternative != nil {
		return Eval(conditionalExp.Alternative, env)
	}

	return NULL
}

func evalIdentifier(ident *ast.Identifier, env *object.Environment) object.Object {
	obj, ok := env.Get(ident.Value)
	if ok {
		return obj
	}

	return newError(IDENT_NOT_FOUND, ident.Value)
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

func evalExpressions(args []ast.Expression, env *object.Environment) []object.Object {
	argsEval := []object.Object{}

	for _, arg := range args {
		argEval := Eval(arg, env)

		if isError(argEval) {
			return []object.Object{argEval}
		}

		argsEval = append(argsEval, argEval)
	}

	return argsEval
}

func applyFunction(fn object.Object, args []object.Object) object.Object {

	function, ok := fn.(*object.Function)
	if !ok {
		return newError(NOT_A_FUNC, fn.Type())
	}

	extendedEnv := extendedFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)

	return unwrapReturnValue(evaluated)
}

func extendedFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for idx, param := range fn.Parameters {
		env.Set(param.Value, args[idx])
	}

	return env
}

func unwrapReturnValue(val object.Object) object.Object {
	returnVal, ok := val.(*object.ReturnValue)
	if ok {
		return returnVal.Value
	}

	return val
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

	// Mismatch case
	if l.Type() != r.Type() {
		return newError(TYPE_MISMATCH_INFIX_MSG, l.Type(), token.PLUS, r.Type())
	}

	// Fetch for the right infix function
	op, ok := object.OBJECT_INFIX_PLUS_FUNCS[l.Type()]
	if !ok {
		return newError(UNKNOWN_OP_INFIX_MSG, l.Type(), token.PLUS, r.Type())
	}

	return op(l, r)
}

func infixMinus(l, r object.Object) object.Object {

	// Mismatch case
	if l.Type() != r.Type() {
		return newError(UNKNOWN_OP_INFIX_MSG, l.Type(), token.MINUS, r.Type())
	}

	// Fetch for the right infix function
	op, ok := object.OBJECT_INFIX_MINUS_FUNCS[l.Type()]
	if !ok {
		return newError(UNKNOWN_OP_INFIX_MSG, l.Type(), token.MINUS, r.Type())
	}

	return op(l, r)
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
