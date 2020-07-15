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
	return fmt.Sprintf("CopiledFunction[%p]", cf)
}
