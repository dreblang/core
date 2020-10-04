package object

import (
	"fmt"

	"github.com/dreblang/core/token"
)

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return IntegerObj }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}
func (i *Integer) String() string { return fmt.Sprintf("%d", i.Value) }

func (obj *Integer) GetMember(name string) Object {
	return newError("No member named [%s]", name)
}

func (obj *Integer) InfixOperation(operator string, other Object) Object {
	switch operator {
	case token.Plus:
		return obj.Add(other)
	case token.Minus:
		return obj.Subtract(other)
	case token.Asterisk:
		return obj.Multiply(other)
	case token.Slash:
		return obj.Divide(other)
	case token.Percent:
		return obj.Modulo(other)
	case token.LessThan:
		return obj.LessThan(other)
	case token.LessOrEqual:
		return obj.LessOrEqual(other)
	case token.GreaterThan:
		return obj.GreaterThan(other)
	case token.GreaterOrEqual:
		return obj.GreaterOrEqual(other)
	case token.Equal:
		return obj.Equals(other)
	case token.NotEqual:
		return obj.NotEquals(other)
	default:
		return newError("%s: %s %s %s", unknownOperatorError, obj.Type(), operator, other.Type())
	}
}

// Arithmetic operations
func (obj *Integer) Add(other Object) Object {
	switch val := other.(type) {
	case *Integer:
		return &Integer{
			Value: obj.Value + val.Value,
		}
	case *Float:
		return &Float{
			Value: float64(obj.Value) + val.Value,
		}
	}

	return newError("%s: %s + %s", typeMissMatchError, obj.Type(), other.Type())
}

func (obj *Integer) Subtract(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return &Integer{
			Value: obj.Value - other.(*Integer).Value,
		}
	case FloatObj:
		return &Float{
			Value: float64(obj.Value) - other.(*Float).Value,
		}
	}
	return newError("%s: %s - %s", typeMissMatchError, obj.Type(), other.Type())
}

func (obj *Integer) Multiply(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return &Integer{
			Value: obj.Value * other.(*Integer).Value,
		}
	case FloatObj:
		return &Float{
			Value: float64(obj.Value) * other.(*Float).Value,
		}
	}
	return newError("%s: %s * %s", typeMissMatchError, obj.Type(), other.Type())
}

func (obj *Integer) Divide(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return &Integer{
			Value: obj.Value / other.(*Integer).Value,
		}
	case FloatObj:
		return &Float{
			Value: float64(obj.Value) / other.(*Float).Value,
		}
	}
	return newError("%s: %s / %s", typeMissMatchError, obj.Type(), other.Type())
}

func (obj *Integer) Modulo(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return &Integer{
			Value: obj.Value % other.(*Integer).Value,
		}
	}
	return newError("%s: %s %% %s", typeMissMatchError, obj.Type(), other.Type())
}

func (obj *Integer) LessThan(other Object) Object {
	switch val := other.(type) {
	case *Integer:
		return NativeBoolToBooleanObject(obj.Value < val.Value)
	case *Float:
		return NativeBoolToBooleanObject(float64(obj.Value) < val.Value)
	}
	return newError("%s: %s < %s", typeMissMatchError, obj.Type(), other.Type())
}

func (obj *Integer) LessOrEqual(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return NativeBoolToBooleanObject(obj.Value <= other.(*Integer).Value)
	case FloatObj:
		return NativeBoolToBooleanObject(float64(obj.Value) <= other.(*Float).Value)
	}
	return newError("%s: %s <= %s", typeMissMatchError, obj.Type(), other.Type())
}

func (obj *Integer) GreaterThan(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return NativeBoolToBooleanObject(obj.Value > other.(*Integer).Value)
	case FloatObj:
		return NativeBoolToBooleanObject(float64(obj.Value) > other.(*Float).Value)
	}
	return newError("%s: %s > %s", typeMissMatchError, obj.Type(), other.Type())
}

func (obj *Integer) GreaterOrEqual(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return NativeBoolToBooleanObject(obj.Value >= other.(*Integer).Value)
	case FloatObj:
		return NativeBoolToBooleanObject(float64(obj.Value) >= other.(*Float).Value)
	}
	return newError("%s: %s >= %s", typeMissMatchError, obj.Type(), other.Type())
}

func (obj *Integer) Equals(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return NativeBoolToBooleanObject(obj.Value == other.(*Integer).Value)
	case FloatObj:
		return NativeBoolToBooleanObject(float64(obj.Value) == other.(*Float).Value)
	}
	return False
}

func (obj *Integer) NotEquals(other Object) Object {
	switch other.Type() {
	case IntegerObj:
		return NativeBoolToBooleanObject(obj.Value != other.(*Integer).Value)
	case FloatObj:
		return NativeBoolToBooleanObject(float64(obj.Value) != other.(*Float).Value)
	}
	return True
}
