package object

import (
	"bytes"
	"strings"
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

func (obj *Array) GetMember(name string) Object {
	switch name {
	case "length":
		return &Integer{Value: int64(len(obj.Elements))}
	}

	return newError("No member named [%s]", name)
}

func (obj *Array) InfixOperation(operator string, other Object) Object {
	return newError("Unsupported operation [%s]", operator)
}
