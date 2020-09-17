package object

import (
	"hash/fnv"

	"github.com/dreblang/core/token"
)

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return StringObj }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
func (s *String) String() string { return s.Value }

func (obj *String) GetMember(name string) Object {
	switch name {
	case "length":
		return &Integer{Value: int64(len(obj.Value))}
	}

	return newError("No member named [%s]", name)
}

func (obj *String) InfixOperation(operator string, other Object) Object {
	switch operator {
	case token.Plus:
		return obj.Add(other)
	}
	return newError("%s: %s %s %s", unknownOperatorError, obj.Type(), operator, other.Type())
}

func (obj *String) Add(other Object) Object {
	switch other.Type() {
	case StringObj:
		return &String{
			Value: obj.Value + other.(*String).Value,
		}
	}
	return newError("Could not concat string with type [%s]", other.Type())
}
