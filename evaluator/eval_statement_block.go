package evaluator

import (
	"github.com/dreblang/core/ast"
	"github.com/dreblang/core/object"
)

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.ReturnValueObj ||
				rt == object.ErrorObj {
				return result
			}
		}
	}

	return result
}
