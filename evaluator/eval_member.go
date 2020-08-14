package evaluator

import (
	"github.com/dreblang/core/ast"
	"github.com/dreblang/core/object"
)

func evalMemberOperation(obj object.Object, member ast.Expression) object.Object {
	// TODO: Implement member execution
	return obj.GetMember(member.String())
}
