package evaluator

import (
	"github.com/dreblang/core/object"
)

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements))

	if idx < -max || idx >= max {
		return Null
	}

	if idx > 0 {
		return arrayObject.Elements[idx]
	}
	return arrayObject.Elements[max+idx]
}
