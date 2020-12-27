package object

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BuiltinObj }
func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) String() string   { return "builtin" }

func (obj *Builtin) GetMember(name string) Object {
	return newError("No member named [%s]", name)
}

func (obj *Builtin) SetMember(name string, value Object) Object {
	return newError("No member named [%s]", name)
}

func (obj *Builtin) Equals(other Object) bool {
	if otherObj, ok := other.(*Builtin); ok {
		return &obj.Fn == &otherObj.Fn
	}
	return false
}
