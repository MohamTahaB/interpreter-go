package object

import "strings"

type InfixFunc[T any] func(a, b T) T

var OBJECT_INFIX_PLUS_FUNCS map[ObjectType]InfixFunc[Object] = map[ObjectType]InfixFunc[Object]{
	INTEGER_OBJ: infixPlusInteger,
	STRING_OBJ:  infixPlusString,
}

var OBJECT_INFIX_MINUS_FUNCS map[ObjectType]InfixFunc[Object] = map[ObjectType]InfixFunc[Object]{
	INTEGER_OBJ: infixMinusInteger,
}

// Define Infix Functions

// Plus

func infixPlusInteger(a, b Object) Object {
	// Cast into an Ingeter
	IntA := a.(*Integer)
	IntB := b.(*Integer)

	return &Integer{
		Value: IntA.Value + IntB.Value,
	}
}

func infixPlusString(a, b Object) Object {
	// Cast into a String
	StrA := a.(*String)
	StrB := b.(*String)

	var buf strings.Builder

	buf.WriteString(StrA.Value)
	buf.WriteString(StrB.Value)

	return &String{
		Value: buf.String(),
	}
}

// Minus

func infixMinusInteger(a, b Object) Object {
	// Cast into an Integer
	IntA := a.(*Integer)
	IntB := b.(*Integer)

	return &Integer{
		Value: IntA.Value - IntB.Value,
	}
}
