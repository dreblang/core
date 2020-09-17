package object

import (
	"fmt"

	"github.com/dreblang/core/token"
)

type Float struct {
	Value float64
}

func (i *Float) Type() ObjectType { return FloatObj }
func (i *Float) Inspect() string  { return fmt.Sprintf("%f", i.Value) }
func (i *Float) String() string   { return fmt.Sprintf("%g", i.Value) }

func (obj *Float) GetMember(name string) Object {
	return newError("No member named [%s]", name)
}

func (obj *Float) InfixOperation(operator string, other Object) Object {
	switch operator {
	case token.Plus:
		return obj.Add(other)
	case token.Minus:
		return obj.Subtract(other)
	case token.Asterisk:
		return obj.Multiply(other)
	case token.Slash:
		return obj.Divide(other)
	case token.LessThan:
		return obj.LessThan(other)
	case token.GreaterThan:
		return obj.GreaterThan(other)
	case token.Equal:
		return obj.Equals(other)
	case token.NotEqual:
		return obj.NotEquals(other)
	default:
		return newError("%s: %s %s %s", unknownOperatorError, obj.Type(), operator, other.Type())
	}
}

// Arithmetic operations
func (obj *Float) Add(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return &Float{
			Value: obj.Value + float64(other.(*Integer).Value),
		}
	case FloatObj:
		return &Float{
			Value: float64(obj.Value) + other.(*Float).Value,
		}
	}
	return newError("Could not perform arithmetic operation")
}

func (obj *Float) Subtract(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return &Float{
			Value: obj.Value - float64(other.(*Integer).Value),
		}
	case FloatObj:
		return &Float{
			Value: float64(obj.Value) - other.(*Float).Value,
		}
	}
	return newError("Could not perform arithmetic operation")
}

func (obj *Float) Multiply(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return &Float{
			Value: obj.Value * float64(other.(*Integer).Value),
		}
	case FloatObj:
		return &Float{
			Value: float64(obj.Value) * other.(*Float).Value,
		}
	}
	return newError("Could not perform arithmetic operation")
}

func (obj *Float) Divide(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return &Float{
			Value: obj.Value / float64(other.(*Integer).Value),
		}
	case FloatObj:
		return &Float{
			Value: float64(obj.Value) / other.(*Float).Value,
		}
	}
	return newError("Could not perform arithmetic operation")
}

func (obj *Float) LessThan(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return NativeBoolToBooleanObject(obj.Value < float64(other.(*Integer).Value))
	case FloatObj:
		return NativeBoolToBooleanObject(obj.Value < other.(*Float).Value)
	}
	return newError("Could not perform arithmetic operation")
}

func (obj *Float) GreaterThan(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return NativeBoolToBooleanObject(obj.Value > float64(other.(*Integer).Value))
	case FloatObj:
		return NativeBoolToBooleanObject(obj.Value > other.(*Float).Value)
	}
	return newError("Could not perform arithmetic operation")
}

func (obj *Float) Equals(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return NativeBoolToBooleanObject(obj.Value == float64(other.(*Integer).Value))
	case FloatObj:
		return NativeBoolToBooleanObject(obj.Value == other.(*Float).Value)
	}
	return newError("Could not perform arithmetic operation")
}

func (obj *Float) NotEquals(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return NativeBoolToBooleanObject(obj.Value != float64(other.(*Integer).Value))
	case FloatObj:
		return NativeBoolToBooleanObject(obj.Value != other.(*Float).Value)
	}
	return newError("Could not perform arithmetic operation")
}
