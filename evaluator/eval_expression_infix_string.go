package evaluator

import (
	"github.com/dreblang/core/object"
	"github.com/dreblang/core/token"
)

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	if operator != token.Plus {
		return newError("%s: %s %s %s", unknownOperatorError, left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ArrayObj && index.Type() == object.IntegerObj:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HashObj:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}
