package object

import (
	"fmt"
)

type Closure struct {
	Fn   *CompiledFunction
	Free []Object
}

func (c *Closure) Type() ObjectType { return ClosureObj }

func (c *Closure) Inspect() string {
	return fmt.Sprintf("Closure[%p]", c)
}

func (c *Closure) String() string {
	return "closure"
}