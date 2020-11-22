package object

import (
	"fmt"

	"github.com/dreblang/core/code"
)

type CompiledFunction struct {
	Instructions  code.Instructions
	NumLocals     int
	NumParameters int
}

func (cf *CompiledFunction) Type() ObjectType { return CompiledFunctionObj }
func (cf *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", cf)
}
func (cf *CompiledFunction) String() string { return "cfunc" }

func (obj *CompiledFunction) GetMember(name string) Object {
	return newError("No member named [%s]", name)
}
func (obj *CompiledFunction) SetMember(name string, value Object) Object {
	return newError("No member named [%s]", name)
}
