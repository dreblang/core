package evaluator

import (
	"github.com/dreblang/core/object"
)

func evalArrayIndexExpression(array, index, indexUpper, indexSkip object.Object, hasUpper, hasSkip bool) object.Object {
	arrayObject := array.(*object.Array)

	var idx int64 = 0
	if index != nil {
		idx = index.(*object.Integer).Value
	}
	max := int64(len(arrayObject.Elements))

	var idxUpper int64 = max

	if idx < 0 {
		idx = max + idx
	}

	if !hasUpper {
		if idx < 0 || idx >= max {
			return Null
		}

		return arrayObject.Elements[idx]
	}

	if indexUpper != nil {
		idxUpper = indexUpper.(*object.Integer).Value
	}
	if idxUpper < 0 {
		idxUpper += max
	}
	var inc int64 = 1
	if hasSkip {
		inc = indexSkip.(*object.Integer).Value
	}

	elements := make([]object.Object, 0)

	for i := idx; i < idxUpper; i += inc {
		elements = append(elements, arrayObject.Elements[i])
	}
	return &object.Array{
		Elements: elements,
	}
}

func evalArrayIndexSetExpression(array, index, indexUpper, indexSkip object.Object, hasUpper, hasSkip bool, value object.Object) object.Object {
	arrayObject := array.(*object.Array)

	var idx int64 = 0
	if index != nil {
		idx = index.(*object.Integer).Value
	}
	max := int64(len(arrayObject.Elements))

	var idxUpper int64 = max

	if idx < 0 {
		idx = max + idx
	}

	if !hasUpper {
		if idx < 0 || idx >= max {
			return Null
		}
		arrayObject.Elements[idx] = value

		return value
	}

	if indexUpper != nil {
		idxUpper = indexUpper.(*object.Integer).Value
	}
	if idxUpper < 0 {
		idxUpper += max
	}
	var inc int64 = 1
	if hasSkip {
		inc = indexSkip.(*object.Integer).Value
	}

	elements := make([]object.Object, 0)

	for i := idx; i < idxUpper; i += inc {
		elements = append(elements, arrayObject.Elements[i])
	}
	return &object.Array{
		Elements: elements,
	}
}
