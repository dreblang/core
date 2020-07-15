package evaluator

import (
	"github.com/dreblang/core/ast"
	"github.com/dreblang/core/object"
)

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("%s: %s", identifierNotFoundError, node.Value)
}
