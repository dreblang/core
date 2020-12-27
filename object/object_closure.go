package object

import (
	"fmt"
)

type Closure struct {
	Fn      *CompiledFunction
	Free    []Object
	Exports map[string]Object
}

func (c *Closure) Type() ObjectType { return ClosureObj }

func (c *Closure) Inspect() string {
	return fmt.Sprintf("Closure[%p]", c)
}

func (c *Closure) String() string {
	return "closure"
}

func (obj *Closure) Equals(other Object) bool {
	if otherObj, ok := other.(*Closure); ok {
		return obj.Fn == otherObj.Fn
	}
	return false
}

func (obj *Closure) GetMember(name string) Object {
	return newError("No member named [%s]", name)
}
func (obj *Closure) SetMember(name string, value Object) Object {
	return newError("No member named [%s]", name)
}
