package object

import "fmt"

type Float struct {
	Value float64
}

func (i *Float) Type() ObjectType { return FloatObj }
func (i *Float) Inspect() string  { return fmt.Sprintf("%f", i.Value) }
func (i *Float) String() string   { return fmt.Sprintf("%g", i.Value) }

func (obj *Float) GetMember(name string) Object {
	return newError("No member named [%s]", name)
}
