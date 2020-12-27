package object

type MemberFunction func(this Object, args ...Object) Object

type MemberFn struct {
	Obj Object
	Fn  MemberFunction
}

func (b *MemberFn) Type() ObjectType { return BuiltinObj }
func (b *MemberFn) Inspect() string  { return "member function" }
func (b *MemberFn) String() string   { return "member" }

func (obj *MemberFn) GetMember(name string) Object {
	return newError("No member named [%s]", name)
}

func (obj *MemberFn) SetMember(name string, value Object) Object {
	return newError("No member named [%s]", name)
}

func (obj *MemberFn) Native() interface{} {
	return nil
}

func (obj *MemberFn) Equals(other Object) bool {
	return false
}
