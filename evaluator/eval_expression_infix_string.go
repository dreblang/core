package evaluator

import (
	"github.com/dreblang/core/object"
)

func evalIndexExpression(left, index, indexUpper, indexSkip object.Object, hasUpper, hasSkip bool) object.Object {
	switch {
	case left.Type() == object.ArrayObj:
		return evalArrayIndexExpression(left, index, indexUpper, indexSkip, hasUpper, hasSkip)

	case left.Type() == object.HashObj:
		return evalHashIndexExpression(left, index)

	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalIndexSetExpression(left, index, indexUpper, indexSkip object.Object, hasUpper, hasSkip bool, value object.Object) object.Object {
	switch {
	case left.Type() == object.ArrayObj:
		return evalArrayIndexSetExpression(left, index, indexUpper, indexSkip, hasUpper, hasSkip, value)

	case left.Type() == object.HashObj:
		return evalHashIndexExpression(left, index)

	default:
		return newError("index operator not supported: %s", left.Type())
	}
}
