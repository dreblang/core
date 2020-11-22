package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dreblang/core/token"
)

type HashKey struct {
	Type  ObjectType
	Value string
}

type Hashable interface {
	HashKey() HashKey
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HashObj }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	var pairs []string
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString(token.LeftBrace)
	out.WriteString(strings.Join(pairs, token.Comma+" "))
	out.WriteString(token.RightBrace)

	return out.String()
}
func (h *Hash) String() string { return "hash" }

func (obj *Hash) GetMember(name string) Object {
	switch name {
	case "length":
		return &Integer{Value: int64(len(obj.Pairs))}
	}

	if val, ok := obj.Pairs[(&String{Value: name}).HashKey()]; ok {
		return val.Value
	}

	return newError("No member named [%s]", name)
}

func (obj *Hash) Native() interface{} {
	result := map[interface{}]interface{}{}
	for _, v := range obj.Pairs {
		result[v.Key.(NativeObject).Native()] = v.Value.(NativeObject).Native()
	}
	return result
}
