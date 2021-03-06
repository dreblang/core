package object

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/dreblang/core/token"
)

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ArrayObj }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	var elements []string
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
func (ao *Array) String() string {
	return "array"
}

func (ao *Array) MarshalJSON() (text []byte, err error) {
	return json.Marshal(ao.Elements)
}

func (obj *Array) GetMember(name string) Object {
	switch name {
	case "length":
		return &Integer{Value: int64(len(obj.Elements))}
	}

	return newError("No member named [%s]", name)
}

func (obj *Array) SetMember(name string, value Object) Object {
	return newError("No member named [%s]", name)
}

func (obj *Array) Native() interface{} {
	result := make([]interface{}, len(obj.Elements))
	for i := range obj.Elements {
		result[i] = obj.Elements[i].(NativeObject).Native()
	}
	return result
}

func (obj *Array) Equals(other Object) bool {
	if otherObj, ok := other.(*Array); ok {
		if len(otherObj.Elements) != len(obj.Elements) {
			return false
		}

		for idx := range obj.Elements {
			if !obj.Elements[idx].Equals(otherObj.Elements[idx]) {
				return false
			}
		}
		return true
	}
	return false
}

func (obj *Array) InfixOperation(operator string, other Object) Object {
	switch operator {
	case token.Plus:
		switch val := other.(type) {
		case *Array:
			return &Array{
				Elements: append(obj.Elements, val.Elements...),
			}
		}
	}

	return newError("%s: %s %s %s", unknownOperatorError, obj.Type(), operator, other.Type())
}
