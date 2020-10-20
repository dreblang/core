package compiler

import (
	"github.com/dreblang/core/object"
)

var coreModules = map[string]func() *object.Scope{}

func RegisterLib(name string, loader func() *object.Scope) {
	coreModules[name] = loader
}
