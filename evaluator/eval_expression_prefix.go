package evaluator

import (
	"github.com/dreblang/core/object"
	"github.com/dreblang/core/token"
)

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case token.Bang:
		return evalBangOperatorExpression(right)
	case token.Minus:
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("%s: %s%s", unknownOperatorError, operator, right.Type())
	}
}
