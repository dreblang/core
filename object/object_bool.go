package object

import (
	"encoding/json"
	"fmt"

	"github.com/dreblang/core/token"
)

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BooleanObj }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value string

	if b.Value {
		value = "t"
	} else {
		value = "f"
	}

	return HashKey{Type: b.Type(), Value: value, Key: b}
}
func (b *Boolean) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}
func (b *Boolean) MarshalJSON() (text []byte, err error) {
	return json.Marshal(b.Value)
}

func (obj *Boolean) GetMember(name string) Object {
	return newError("No member named [%s]", name)
}

func (obj *Boolean) SetMember(name string, value Object) Object {
	return newError("No member named [%s]", name)
}

func (obj *Boolean) Equals(other Object) bool {
	if otherObj, ok := other.(*Boolean); ok {
		return obj.Value == otherObj.Value
	}
	return false
}

func (obj *Boolean) Native() interface{} {
	return obj.Value
}

func (obj *Boolean) InfixOperation(operator string, other Object) Object {
	switch operator {
	case token.Equal:
		switch val := other.(type) {
		case *Boolean:
			return NativeBoolToBooleanObject(obj.Value == val.Value)
		default:
			return False
		}

	case token.NotEqual:
		switch val := other.(type) {
		case *Boolean:
			return NativeBoolToBooleanObject(obj.Value != val.Value)
		default:
			return True
		}
	}
	return newError("%s: %s %s %s", unknownOperatorError, obj.Type(), operator, other.Type())
}
