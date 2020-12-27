package object

import (
	"bytes"
	"fmt"

	"github.com/dreblang/core/token"
)

type Bytes struct {
	Value []byte
}

func (s *Bytes) Type() ObjectType { return BytesObj }
func (s *Bytes) Inspect() string  { return fmt.Sprintf("bytes(%s)", string(s.Value)) }
func (s *Bytes) String() string   { return string(s.Value) }

func (obj *Bytes) GetMember(name string) Object {
	switch name {
	case "length":
		return &Integer{Value: int64(len(obj.Value))}

	case "sub":
		return &MemberFn{
			Obj: obj,
			Fn:  bytesSub,
		}

	case "starts_with":
		return &MemberFn{
			Obj: obj,
			Fn:  bytesStartsWith,
		}

	case "ends_with":
		return &MemberFn{
			Obj: obj,
			Fn:  bytesEndsWith,
		}
	}

	return newError("No member named [%s]", name)
}

func (obj *Bytes) SetMember(name string, value Object) Object {
	return newError("No member named [%s]", name)
}

func (obj *Bytes) Native() interface{} {
	return obj.Value
}

func (obj *Bytes) Equals(other Object) bool {
	if otherObj, ok := other.(*Bytes); ok {
		return bytes.Compare(obj.Value, otherObj.Value) == 0
	}
	return false
}

func (obj *Bytes) InfixOperation(operator string, other Object) Object {
	switch operator {
	case token.Plus:
		switch val := other.(type) {
		case *Bytes:
			return &Bytes{
				Value: append(obj.Value, val.Value...),
			}
		}

	case token.LessThan:
		switch val := other.(type) {
		case *Bytes:
			return NativeBoolToBooleanObject(bytes.Compare(obj.Value, val.Value) < 0)
		}

	case token.LessOrEqual:
		switch val := other.(type) {
		case *Bytes:
			return NativeBoolToBooleanObject(bytes.Compare(obj.Value, val.Value) <= 0)
		}

	case token.GreaterThan:
		switch val := other.(type) {
		case *Bytes:
			return NativeBoolToBooleanObject(bytes.Compare(obj.Value, val.Value) > 0)
		}

	case token.GreaterOrEqual:
		switch val := other.(type) {
		case *Bytes:
			return NativeBoolToBooleanObject(bytes.Compare(obj.Value, val.Value) >= 0)
		}

	case token.Equal:
		switch val := other.(type) {
		case *Bytes:
			return NativeBoolToBooleanObject(bytes.Compare(obj.Value, val.Value) == 0)
		default:
			return False
		}

	}
	return newError("%s: %s %s %s", unknownOperatorError, obj.Type(), operator, other.Type())
}

func bytesSub(this Object, args ...Object) Object {
	str := this.(*Bytes)
	switch len(args) {
	case 0:
		return this
	case 1:
		if startIdx, ok := args[0].(*Integer); ok {
			return &Bytes{
				Value: str.Value[startIdx.Value:],
			}
		}
	case 2:
		if startIdx, ok := args[0].(*Integer); ok {
			if endIdx, ok := args[1].(*Integer); ok {
				return &Bytes{
					Value: str.Value[startIdx.Value:endIdx.Value],
				}
			}
		}
	}
	return newError("Could not execute sub-string operation. Invalid arguments!")
}

func bytesStartsWith(this Object, args ...Object) Object {
	str := this.(*Bytes)
	switch len(args) {
	case 1:
		other := args[0].(*Bytes)
		if len(other.Value) > len(str.Value) {
			return False
		}

		return NativeBoolToBooleanObject(
			bytes.Compare(str.Value[:len(other.Value)], other.Value) == 0,
		)
	}

	return newError("Invalid arguments!")
}

func bytesEndsWith(this Object, args ...Object) Object {
	str := this.(*Bytes)
	switch len(args) {
	case 1:
		other := args[0].(*Bytes)
		if len(other.Value) > len(str.Value) {
			return False
		}

		return NativeBoolToBooleanObject(
			bytes.Compare(str.Value[len(str.Value)-len(other.Value):], other.Value) == 0,
		)
	}

	return newError("Invalid arguments!")
}
