package object

import "github.com/dreblang/core/token"

type Null struct{}

func (n *Null) Type() ObjectType { return NullObj }
func (n *Null) Inspect() string  { return "null" }
func (n *Null) String() string   { return "null" }

var NullObject = &Null{}

func (obj *Null) GetMember(name string) Object {
	return newError("No member named [%s]", name)
}
func (obj *Null) SetMember(name string, value Object) Object {
	return newError("No member named [%s]", name)
}

func (obj *Null) Native() interface{} {
	return nil
}

func (obj *Null) Equals(other Object) bool {
	if _, ok := other.(*Null); ok {
		return true
	}
	return false
}

func (obj *Null) InfixOperation(operator string, other Object) Object {
	switch operator {
	case token.Equal:
		switch other.(type) {
		case *Null:
			return True
		default:
			return False
		}

	case token.NotEqual:
		switch other.(type) {
		case *Null:
			return False
		default:
			return True
		}
	}
	return newError("%s: %s %s %s", unknownOperatorError, obj.Type(), operator, other.Type())
}
