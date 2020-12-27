package object

import (
	"strings"

	"github.com/dreblang/core/token"
)

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return StringObj }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	return HashKey{Type: s.Type(), Value: s.Value}
}
func (s *String) String() string { return s.Value }

func (obj *String) GetMember(name string) Object {
	switch name {
	case "length":
		return &Integer{Value: int64(len(obj.Value))}

	case "sub":
		return &MemberFn{
			Obj: obj,
			Fn:  stringSub,
		}

	case "upper":
		return &MemberFn{
			Obj: obj,
			Fn:  stringUpper,
		}

	case "lower":
		return &MemberFn{
			Obj: obj,
			Fn:  stringLower,
		}

	case "replace":
		return &MemberFn{
			Obj: obj,
			Fn:  stringReplace,
		}

	case "strip":
		return &MemberFn{
			Obj: obj,
			Fn:  stringStrip,
		}

	case "split":
		return &MemberFn{
			Obj: obj,
			Fn:  stringSplit,
		}

	case "starts_with":
		return &MemberFn{
			Obj: obj,
			Fn:  stringStartsWith,
		}

	case "ends_with":
		return &MemberFn{
			Obj: obj,
			Fn:  stringEndsWith,
		}
	}

	return newError("No member named [%s]", name)
}

func (obj *String) SetMember(name string, value Object) Object {
	return newError("No member named [%s]", name)
}

func (obj *String) Native() interface{} {
	return obj.Value
}

func (obj *String) Equals(other Object) bool {
	if otherObj, ok := other.(*String); ok {
		return obj.Value == otherObj.Value
	}
	return false
}

func (obj *String) InfixOperation(operator string, other Object) Object {
	switch operator {
	case token.Plus:
		switch val := other.(type) {
		case *String:
			return &String{
				Value: obj.Value + val.Value,
			}
		}

	case token.LessThan:
		switch val := other.(type) {
		case *String:
			return NativeBoolToBooleanObject(obj.Value < val.Value)
		}

	case token.LessOrEqual:
		switch val := other.(type) {
		case *String:
			return NativeBoolToBooleanObject(obj.Value <= val.Value)
		}

	case token.GreaterThan:
		switch val := other.(type) {
		case *String:
			return NativeBoolToBooleanObject(obj.Value > val.Value)
		}

	case token.GreaterOrEqual:
		switch val := other.(type) {
		case *String:
			return NativeBoolToBooleanObject(obj.Value >= val.Value)
		}

	case token.Equal:
		switch val := other.(type) {
		case *String:
			return NativeBoolToBooleanObject(obj.Value == val.Value)
		default:
			return False
		}

	case token.String:
		switch val := other.(type) {
		case *String:
			return NativeBoolToBooleanObject(obj.Value != val.Value)
		default:
			return True
		}
	}
	return newError("%s: %s %s %s", unknownOperatorError, obj.Type(), operator, other.Type())
}

func stringSub(this Object, args ...Object) Object {
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

func stringUpper(this Object, args ...Object) Object {
	str := this.(*String)
	switch len(args) {
	case 0:
		return &String{
			Value: strings.ToUpper(str.Value),
		}
	}
	return newError("Could not execute string upper operation. Invalid arguments!")
}

func stringLower(this Object, args ...Object) Object {
	str := this.(*String)
	switch len(args) {
	case 0:
		return &String{
			Value: strings.ToLower(str.Value),
		}
	}
	return newError("Could not execute string lower operation. Invalid arguments!")
}

func stringReplace(this Object, args ...Object) Object {
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

func stringStrip(this Object, args ...Object) Object {
	str := this.(*String)
	switch len(args) {
	case 0:
		return &String{
			Value: strings.Trim(str.Value, " \t\n"),
		}
	}
	return newError("Invalid arguments!")
}

func stringSplit(this Object, args ...Object) Object {
	str := this.(*String)
	var values []string
	switch len(args) {
	case 0:
		values = strings.Split(str.Value, " ")
	case 1:
		values = strings.Split(str.Value, args[0].(*String).Value)
	default:
		return newError("Invalid arguments!")
	}

	valueObjs := make([]Object, len(values))
	for i, v := range values {
		valueObjs[i] = &String{Value: v}
	}

	return &Array{
		Elements: valueObjs,
	}
}

func stringStartsWith(this Object, args ...Object) Object {
	str := this.(*String)
	switch len(args) {
	case 1:
		other := args[0].(*String)
		if len(other.Value) > len(str.Value) {
			return False
		}

		return NativeBoolToBooleanObject(
			str.Value[:len(other.Value)] == other.Value,
		)
	}

	return newError("Invalid arguments!")
}

func stringEndsWith(this Object, args ...Object) Object {
	str := this.(*String)
	switch len(args) {
	case 1:
		other := args[0].(*String)
		if len(other.Value) > len(str.Value) {
			return False
		}

		return NativeBoolToBooleanObject(
			str.Value[len(str.Value)-len(other.Value):] == other.Value,
		)
	}

	return newError("Invalid arguments!")
}
