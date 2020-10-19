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

func (obj *Builtin) Native() interface{} {
	return obj.Fn
}
func (obj *Builtin) InfixOperation(operator string, other Object) Object {
	return newError("%s: %s %s %s", unknownOperatorError, obj.Type(), operator, other.Type())
}
