package http

import (
	"github.com/dreblang/core/corelib/http/objects"
	"github.com/dreblang/core/object"
)

var httpScope = &object.Scope{
	Name: "http",
	Exports: map[string]object.Object{
		"createServer": &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				return objects.NewServer()
			},
		},
	},
}

// Load the scope
func Load() *object.Scope {
	return httpScope
}
