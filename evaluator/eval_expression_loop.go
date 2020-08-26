package evaluator

import (
	"github.com/dreblang/core/ast"
	"github.com/dreblang/core/object"
)

func evalLoopExpression(le *ast.LoopExpression, env *object.Environment) object.Object {
	for {
		condition := Eval(le.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			Eval(le.Consequence, env)
		} else {
			return Null
		}
	}
}
