package evaluator

import (
	"github.com/dreblang/core/object"
	"github.com/dreblang/core/token"
)

func evalNumericInfixExpression(operator string, left, right object.Object) object.Object {
	switch left.(type) {
	case *object.Integer:
		switch right.(type) {
		case *object.Integer:
			return evalIntIntExpression(operator, left.(*object.Integer), right.(*object.Integer))
		case *object.Float:
			return evalIntFloatExpression(operator, left.(*object.Integer), right.(*object.Float))
		}

	case *object.Float:
		switch right.(type) {
		case *object.Integer:
			return evalFloatIntExpression(operator, left.(*object.Float), right.(*object.Integer))
		case *object.Float:
			return evalFloatFloatExpression(operator, left.(*object.Float), right.(*object.Float))
		}
	}

	return newError("Unable to perform operation %s between types %T and %T!", operator, left, right)
}

func evalIntIntExpression(operator string, left *object.Integer, right *object.Integer) object.Object {
	leftValue := left.Value
	rightValue := right.Value

	switch operator {
	case token.Plus:
		return &object.Integer{Value: leftValue + rightValue}
	case token.Minus:
		return &object.Integer{Value: leftValue - rightValue}
	case token.Asterisk:
		return &object.Integer{Value: leftValue * rightValue}
	case token.Slash:
		return &object.Integer{Value: leftValue / rightValue}
	case token.LessThan:
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case token.GreaterThan:
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case token.Equal:
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case token.NotEqual:
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("%s: %s %s %s", unknownOperatorError, left.Type(), operator, right.Type())
	}
}

func evalIntFloatExpression(operator string, left *object.Integer, right *object.Float) object.Object {
	leftValue := float64(left.Value)
	rightValue := right.Value

	switch operator {
	case token.Plus:
		return &object.Float{Value: leftValue + rightValue}
	case token.Minus:
		return &object.Float{Value: leftValue - rightValue}
	case token.Asterisk:
		return &object.Float{Value: leftValue * rightValue}
	case token.Slash:
		return &object.Float{Value: leftValue / rightValue}
	case token.LessThan:
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case token.GreaterThan:
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case token.Equal:
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case token.NotEqual:
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("%s: %s %s %s", unknownOperatorError, left.Type(), operator, right.Type())
	}
}

func evalFloatIntExpression(operator string, left *object.Float, right *object.Integer) object.Object {
	leftValue := left.Value
	rightValue := float64(right.Value)

	switch operator {
	case token.Plus:
		return &object.Float{Value: leftValue + rightValue}
	case token.Minus:
		return &object.Float{Value: leftValue - rightValue}
	case token.Asterisk:
		return &object.Float{Value: leftValue * rightValue}
	case token.Slash:
		return &object.Float{Value: leftValue / rightValue}
	case token.LessThan:
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case token.GreaterThan:
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case token.Equal:
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case token.NotEqual:
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("%s: %s %s %s", unknownOperatorError, left.Type(), operator, right.Type())
	}
}

func evalFloatFloatExpression(operator string, left *object.Float, right *object.Float) object.Object {
	leftValue := left.Value
	rightValue := right.Value

	switch operator {
	case token.Plus:
		return &object.Float{Value: leftValue + rightValue}
	case token.Minus:
		return &object.Float{Value: leftValue - rightValue}
	case token.Asterisk:
		return &object.Float{Value: leftValue * rightValue}
	case token.Slash:
		return &object.Float{Value: leftValue / rightValue}
	case token.LessThan:
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case token.GreaterThan:
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case token.Equal:
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case token.NotEqual:
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("%s: %s %s %s", unknownOperatorError, left.Type(), operator, right.Type())
	}
}
