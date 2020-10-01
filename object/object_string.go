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
	case "sub":
		return &MemberFn{
			Obj: obj,
			Fn:  StringSub,
		}
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

func StringSub(this Object, args ...Object) Object {
	str := this.(*String)
	switch len(args) {
	case 0:
		return this
	case 1:
		if startIdx, ok := args[0].(*Integer); ok {
			return &String{
				Value: str.Value[startIdx.Value:],
			}
		}
	case 2:
		if startIdx, ok := args[0].(*Integer); ok {
			if endIdx, ok := args[1].(*Integer); ok {
				return &String{
					Value: str.Value[startIdx.Value:endIdx.Value],
				}
			}
		}
	}
	return newError("Could not execute sub string operation. Invalid arguments!")
}
