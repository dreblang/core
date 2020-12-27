package object

import (
	"fmt"

	"github.com/dreblang/core/code"
)

type Scope struct {
	Name         string
	Constants    []Object
	Instructions code.Instructions
	NumLocals    int
	Exports      map[string]Object
}

func (cf *Scope) Type() ObjectType { return ScopeObj }
func (cf *Scope) Inspect() string {
	return fmt.Sprintf("Scope[%s]", cf.Name)
}
func (cf *Scope) String() string { return "scope" }

func (obj *Scope) GetMember(name string) Object {
	if res, ok := obj.Exports[name]; ok {
		return res
	}
	return newError("No member named [%s]", name)
}
func (obj *Scope) SetMember(name string, value Object) Object {
	return newError("No member named [%s]", name)
}

func (obj *Scope) Native() interface{} {
	return nil
}

func (obj *Scope) InfixOperation(operator string, other Object) Object {
	return newError("%s: %s %s %s", unknownOperatorError, obj.Type(), operator, other.Type())
}

func (obj *Scope) Equals(other Object) bool {
	return false
}
