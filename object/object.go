package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"

	RETURN_OBJ = "RETURN_VAL"
	ERROR_OBJ  = "ERROR"
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
	}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}
