package object

import (
	"encoding/json"
	"fmt"

	"github.com/dreblang/core/token"
)

type Float struct {
	Value float64
}

func (i *Float) Type() ObjectType { return FloatObj }
func (i *Float) Inspect() string  { return fmt.Sprintf("%f", i.Value) }
func (i *Float) String() string   { return fmt.Sprintf("%g", i.Value) }
func (i *Float) MarshalText() (text []byte, err error) {
	return json.Marshal(i.Value)
}

func (i *Float) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: fmt.Sprint(i.Value)}
}

func (obj *Float) GetMember(name string) Object {
	return newError("No member named [%s]", name)
}
func (obj *Float) SetMember(name string, value Object) Object {
	return newError("No member named [%s]", name)
}

func (obj *Float) Native() interface{} {
	return obj.Value
}

func (obj *Float) Equals(other Object) bool {
	if otherObj, ok := other.(*Float); ok {
		return obj.Value == otherObj.Value
	}

	if otherObj, ok := other.(*Integer); ok {
		return obj.Value == float64(otherObj.Value)
	}

	return false
}

func (obj *Float) InfixOperation(operator string, other Object) Object {
	switch operator {
	case token.Plus:
		switch val := other.(type) {
		case *Integer:
			return &Float{
				Value: obj.Value + float64(val.Value),
			}
		case *Float:
			return &Float{
				Value: obj.Value + val.Value,
			}
		}

	case token.Minus:
		switch val := other.(type) {
		case *Integer:
			return &Float{
				Value: obj.Value - float64(val.Value),
			}
		case *Float:
			return &Float{
				Value: obj.Value - val.Value,
			}
		}

	case token.Asterisk:
		switch val := other.(type) {
		case *Integer:
			return &Float{
				Value: obj.Value * float64(val.Value),
			}
		case *Float:
			return &Float{
				Value: obj.Value * val.Value,
			}
		}

	case token.Slash:
		switch val := other.(type) {
		case *Integer:
			return &Float{
				Value: obj.Value / float64(val.Value),
			}
		case *Float:
			return &Float{
				Value: obj.Value / val.Value,
			}
		}

	case token.LessThan:
		switch val := other.(type) {
		case *Integer:
			return NativeBoolToBooleanObject(obj.Value < float64(val.Value))
		case *Float:
			return NativeBoolToBooleanObject(obj.Value < val.Value)
		}

	case token.LessOrEqual:
		switch val := other.(type) {
		case *Integer:
			return NativeBoolToBooleanObject(obj.Value <= float64(val.Value))
		case *Float:
			return NativeBoolToBooleanObject(obj.Value <= val.Value)
		}

	case token.GreaterThan:
		switch val := other.(type) {
		case *Integer:
			return NativeBoolToBooleanObject(obj.Value > float64(val.Value))
		case *Float:
			return NativeBoolToBooleanObject(obj.Value > val.Value)
		}

	case token.GreaterOrEqual:
		switch val := other.(type) {
		case *Integer:
			return NativeBoolToBooleanObject(obj.Value >= float64(val.Value))
		case *Float:
			return NativeBoolToBooleanObject(obj.Value >= val.Value)
		}

	case token.Equal:
		switch val := other.(type) {
		case *Integer:
			return NativeBoolToBooleanObject(obj.Value == float64(val.Value))
		case *Float:
			return NativeBoolToBooleanObject(obj.Value == val.Value)
		default:
			return False
		}

	case token.NotEqual:
		switch val := other.(type) {
		case *Integer:
			return NativeBoolToBooleanObject(obj.Value != float64(val.Value))
		case *Float:
			return NativeBoolToBooleanObject(obj.Value != val.Value)
		default:
			return True
		}
	}
	return newError("%s: %s %s %s", typeMissMatchError, obj.Type(), operator, other.Type())
}
