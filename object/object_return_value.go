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
