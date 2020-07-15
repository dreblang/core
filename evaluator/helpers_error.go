package evaluator

import (
	"fmt"

	"github.com/dreblang/core/object"
)

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorObj
	} else {
		return false
	}
}
