package evaluator

import (
	"github.com/dreblang/core/object"
	"github.com/dreblang/core/token"
)

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case (left.Type() == object.IntegerObj || left.Type() == object.FloatObj) && (right.Type() == object.IntegerObj || right.Type() == object.FloatObj):
		return evalNumericInfixExpression(operator, left, right)

	case left.Type() == object.StringObj && right.Type() == object.StringObj:
		return evalStringInfixExpression(operator, left, right)

	case operator == token.Equal:
		return nativeBoolToBooleanObject(left == right)

	case operator == token.NotEqual:
		return nativeBoolToBooleanObject(left != right)

	case left.Type() != right.Type():
		return newError("%s: %s %s %s", typeMissMatchError, left.Type(), operator, right.Type())

	default:
		return newError("%s: %s %s %s", unknownOperatorError, left.Type(), operator, right.Type())
	}
}
