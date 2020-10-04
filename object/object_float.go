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
