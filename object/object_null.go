package object

type Null struct{}

func (n *Null) Type() ObjectType { return NullObj }
func (n *Null) Inspect() string  { return "null" }
func (n *Null) String() string   { return "null" }

var NullObject = &Null{}

func (obj *Null) GetMember(name string) Object {
	return newError("No member named [%s]", name)
}
