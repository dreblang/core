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
	return HashKey{Type: i.Type(), Value: fmt.Sprint(i.Value)}
}
func (i *Integer) String() string { return fmt.Sprintf("%d", i.Value) }

func (obj *Integer) GetMember(name string) Object {
	return newError("No member named [%s]", name)
}

func (obj *Integer) SetMember(name string, value Object) Object {
	return newError("No member named [%s]", name)
}

func (obj *Integer) Native() interface{} {
	return obj.Value
}

func (obj *Integer) Equals(other Object) bool {
	if otherObj, ok := other.(*Integer); ok {
		return obj.Value == otherObj.Value
	}

	if otherObj, ok := other.(*Float); ok {
		return float64(obj.Value) == otherObj.Value
	}
	return false
}

func (obj *Integer) InfixOperation(operator string, other Object) Object {
	switch operator {
	case token.Plus:
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

	case token.Minus:
		switch val := other.(type) {
		case *Integer:
			return &Integer{
				Value: obj.Value - val.Value,
			}
		case *Float:
			return &Float{
				Value: float64(obj.Value) - val.Value,
			}
		}

	case token.Asterisk:
		switch val := other.(type) {
		case *Integer:
			return &Integer{
				Value: obj.Value * val.Value,
			}
		case *Float:
			return &Float{
				Value: float64(obj.Value) * val.Value,
			}
		}

	case token.Slash:
		switch val := other.(type) {
		case *Integer:
			return &Integer{
				Value: obj.Value / val.Value,
			}
		case *Float:
			return &Float{
				Value: float64(obj.Value) / val.Value,
			}
		}

	case token.Percent:
		switch val := other.(type) {
		case *Integer:
			return &Integer{
				Value: obj.Value % val.Value,
			}
		}

	case token.LessThan:
		switch val := other.(type) {
		case *Integer:
			return NativeBoolToBooleanObject(obj.Value < val.Value)
		case *Float:
			return NativeBoolToBooleanObject(float64(obj.Value) < val.Value)
		}

	case token.LessOrEqual:
		switch val := other.(type) {
		case *Integer:
			return NativeBoolToBooleanObject(obj.Value <= val.Value)
		case *Float:
			return NativeBoolToBooleanObject(float64(obj.Value) <= val.Value)
		}

	case token.GreaterThan:
		switch val := other.(type) {
		case *Integer:
			return NativeBoolToBooleanObject(obj.Value > val.Value)
		case *Float:
			return NativeBoolToBooleanObject(float64(obj.Value) > val.Value)
		}

	case token.GreaterOrEqual:
		switch val := other.(type) {
		case *Integer:
			return NativeBoolToBooleanObject(obj.Value >= val.Value)
		case *Float:
			return NativeBoolToBooleanObject(float64(obj.Value) >= val.Value)
		}

	case token.Equal:
		switch val := other.(type) {
		case *Integer:
			return NativeBoolToBooleanObject(obj.Value == val.Value)
		case *Float:
			return NativeBoolToBooleanObject(float64(obj.Value) == val.Value)
		default:
			return False
		}

	case token.NotEqual:
		switch val := other.(type) {
		case *Integer:
			return NativeBoolToBooleanObject(obj.Value != val.Value)
		case *Float:
			return NativeBoolToBooleanObject(float64(obj.Value) != val.Value)
		default:
			return True
		}
	}
	return newError("%s: %s %s %s", typeMissMatchError, obj.Type(), operator, other.Type())
}
