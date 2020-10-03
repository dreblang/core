package object

import (
	"hash/fnv"
	"strings"

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
	case "upper":
		return &MemberFn{
			Obj: obj,
			Fn:  StringUpper,
		}
	case "lower":
		return &MemberFn{
			Obj: obj,
			Fn:  StringLower,
		}
	case "replace":
		return &MemberFn{
			Obj: obj,
			Fn:  StringReplace,
		}
	}

	return newError("No member named [%s]", name)
}

func (obj *String) InfixOperation(operator string, other Object) Object {
	switch operator {
	case token.Plus:
		return obj.Add(other)
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

func (obj *String) LessThan(other Object) Object {
	switch other.Type() {
	case StringObj:
		return NativeBoolToBooleanObject(obj.Value < other.(*String).Value)
	}
	return newError("%s: %s < %s", typeMissMatchError, obj.Type(), other.Type())
}

func (obj *String) LessOrEqual(other Object) Object {
	switch other.Type() {
	case StringObj:
		return NativeBoolToBooleanObject(obj.Value <= other.(*String).Value)
	}
	return newError("%s: %s <= %s", typeMissMatchError, obj.Type(), other.Type())
}

func (obj *String) GreaterThan(other Object) Object {
	switch other.Type() {
	case StringObj:
		return NativeBoolToBooleanObject(obj.Value > other.(*String).Value)
	}
	return newError("%s: %s > %s", typeMissMatchError, obj.Type(), other.Type())
}

func (obj *String) GreaterOrEqual(other Object) Object {
	switch other.Type() {
	case StringObj:
		return NativeBoolToBooleanObject(obj.Value >= other.(*String).Value)
	}
	return newError("%s: %s >= %s", typeMissMatchError, obj.Type(), other.Type())
}

func (obj *String) Equals(other Object) Object {
	switch other.Type() {
	case StringObj:
		return NativeBoolToBooleanObject(obj.Value == other.(*String).Value)
	}
	return False
}

func (obj *String) NotEquals(other Object) Object {
	switch other.Type() {
	case StringObj:
		return NativeBoolToBooleanObject(obj.Value != other.(*String).Value)
	}
	return True
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
	return newError("Could not execute sub-string operation. Invalid arguments!")
}

func StringUpper(this Object, args ...Object) Object {
	str := this.(*String)
	switch len(args) {
	case 0:
		return &String{
			Value: strings.ToUpper(str.Value),
		}
	}
	return newError("Could not execute string upper operation. Invalid arguments!")
}

func StringLower(this Object, args ...Object) Object {
	str := this.(*String)
	switch len(args) {
	case 0:
		return &String{
			Value: strings.ToLower(str.Value),
		}
	}
	return newError("Could not execute string lower operation. Invalid arguments!")
}

func StringReplace(this Object, args ...Object) Object {
	str := this.(*String)
	switch len(args) {
	case 2:
		if search, ok := args[0].(*String); ok {
			if replace, ok := args[1].(*String); ok {
				return &String{
					Value: strings.ReplaceAll(str.Value, search.Value, replace.Value),
				}
			}
		}
	}
	return newError("Could not execute string replace operation. Invalid arguments!")
}
