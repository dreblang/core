package evaluator

import (
	"github.com/dreblang/core/object"
)

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	return left.InfixOperation(operator, right)
}
