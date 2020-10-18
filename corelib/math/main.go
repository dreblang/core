package main

import (
	"github.com/dreblang/core/object"

	"math"
)

var mathScope = &object.Scope{
	Name: "math",
	Exports: map[string]object.Object{
		"PI": &object.Float{Value: math.Pi},
		"sin": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if x, ok := args[0].(*object.Float); ok {
					return &object.Float{Value: math.Sin(x.Value)}
				} else if x, ok := args[0].(*object.Integer); ok {
					return &object.Float{Value: math.Sin(float64(x.Value))}
				}
				return object.NullObject
			},
		},
	},
}

// Load the scope
func Load() *object.Scope {
	return mathScope
}
