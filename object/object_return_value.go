package object

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return ReturnValueObj }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) String() string   { return rv.Value.String() }

func (obj *ReturnValue) GetMember(name string) Object {
	return newError("No member named [%s]", name)
}

func (obj *ReturnValue) SetMember(name string, value Object) Object {
	return newError("No member named [%s]", name)
}

func (obj *ReturnValue) InfixOperation(operator string, other Object) Object {
	return newError("%s: %s %s %s", unknownOperatorError, obj.Type(), operator, other.Type())
}

func (obj *ReturnValue) Equals(other Object) bool {
	return obj.Value.Equals(other)
}
