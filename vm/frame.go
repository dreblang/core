package vm

import (
	"github.com/dreblang/core/code"
	"github.com/dreblang/core/object"
)

type Frame struct {
	cl          *object.Closure
	ip          int
	basePointer int

	instructions code.Instructions
}

func NewFrame(cl *object.Closure, basePointer int) *Frame {
	return &Frame{
		cl:          cl,
		ip:          -1,
		basePointer: basePointer,

		instructions: cl.Fn.Instructions,
	}
}
