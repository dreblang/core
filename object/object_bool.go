package object

import (
	"fmt"

	"github.com/dreblang/core/token"
)

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BooleanObj }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}
func (b *Boolean) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

func (obj *Boolean) GetMember(name string) Object {
	return newError("No member named [%s]", name)
}

func (obj *Boolean) InfixOperation(operator string, other Object) Object {
	switch operator {
	case token.Equal:
		return obj.Equals(other)
	case token.NotEqual:
		return obj.NotEquals(other)
	}
	return newError("%s: %s %s %s", unknownOperatorError, obj.Type(), operator, other.Type())
}

func (obj *Boolean) Equals(other Object) Object {
	switch other.Type() {
	case BooleanObj:
		return NativeBoolToBooleanObject(obj.Value == other.(*Boolean).Value)
	}
	return newError("Could not perform arithmetic operation")
}

func (obj *Boolean) NotEquals(other Object) Object {
	switch other.Type() {
	case BooleanObj:
		return NativeBoolToBooleanObject(obj.Value != other.(*Boolean).Value)
	}
	return newError("Could not perform arithmetic operation")
}
