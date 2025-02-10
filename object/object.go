package object

import (
	"fmt"
	"strings"

	"github.com/MohamTahaB/interpreter-go/ast"
)

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"

	RETURN_OBJ   = "RETURN_VAL"
	ERROR_OBJ    = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	Truthy() bool
}

// Integer type
type Integer struct {
	Value int64
}

// Boolean type
type Boolean struct {
	Value bool
}

// Null type
type Null struct{}

// Return wrapper
type ReturnValue struct {
	Value Object
}

// Internal Error Wrapper
type Error struct {
	Message string
}

// Environment
type Environment struct {
	store map[string]Object
	outer *Environment
}

// Function Object
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) Truthy() bool {
	return i.Value != 0
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) Truthy() bool {
	return b.Value
}

func (n *Null) Inspect() string {
	return "null"
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

func (n *Null) Truthy() bool {
	return false
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_OBJ
}

func (rv *ReturnValue) Truthy() bool {
	return rv.Value.Truthy()
}

func (e *Error) Inspect() string {
	return e.Message
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Truthy() bool {
	return false
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: nil,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	if e == nil {
		return nil, false
	}

	obj, ok := e.store[name]
	if ok {
		return obj, ok
	}

	return e.outer.Get(name)
}

func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) Inspect() string {
	var out strings.Builder
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

func (f *Function) Truthy() bool {
	return true
}
