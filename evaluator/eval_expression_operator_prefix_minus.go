package evaluator

import (
	"github.com/dreblang/core/object"
)

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() == object.IntegerObj {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	}

	if right.Type() == object.FloatObj {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}

	return newError("%s: -%s", unknownOperatorError, right.Type())
}
