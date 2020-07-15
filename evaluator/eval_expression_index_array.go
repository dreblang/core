package evaluator

import (
	"github.com/dreblang/core/object"
)

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return Null
	}

	return arrayObject.Elements[idx]
}
