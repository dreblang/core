package object

import (
	"hash/fnv"
)

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return StringObj }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}
func (s *String) String() string { return s.Value }

func (obj *String) GetMember(name string) Object {
	switch name {
	case "length":
		return &Integer{Value: int64(len(obj.Value))}
	}

	return newError("No member named [%s]", name)
}
