package object

import (
	"bytes"
	"strings"

	"github.com/dreblang/core/ast"
)

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FunctionObj }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
func (f *Function) String() string {
	return "func"
}

func (obj *Function) GetMember(name string) Object {
	return newError("No member named [%s]", name)
}

func (obj *Function) InfixOperation(operator string, other Object) Object {
	return newError("%s: %s %s %s", unknownOperatorError, obj.Type(), operator, other.Type())
}
