package object

import (
	"bytes"
	"encoding/json"
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

func (p HashPair) MarshalJSON() (text []byte, err error) {
	buf := bytes.NewBuffer([]byte{})

	b, _ := json.Marshal(p.Key)
	buf.Write(b)

	buf.Write([]byte(":"))

	b, _ = json.Marshal(p.Value)
	buf.Write(b)

	return buf.Bytes(), nil
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

func (h *Hash) MarshalJSON() (text []byte, err error) {
	buf := bytes.NewBuffer([]byte{})
	buf.Write([]byte("{"))

	for _, v := range h.Pairs {
		b, _ := v.MarshalJSON()
		buf.Write(b)
	}

	b := buf.Bytes()
	if len(h.Pairs) > 0 {
		b[len(b)-1] = '}'
	} else {
		b = append(b, '}')
	}

	return b, nil
}

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

func (obj *Hash) SetMember(name string, value Object) Object {
	key := &String{Value: name}
	obj.Pairs[key.HashKey()] = HashPair{
		Key:   key,
		Value: value,
	}
	return value
}

func (obj *Hash) Equals(other Object) bool {
	// if otherObj, ok := other.(*Hash); ok {
	// 	otherObj.Pairs
	// }
	// FIXME: Compare hash
	return false
}

func (obj *Hash) Native() interface{} {
	result := map[interface{}]interface{}{}
	for _, v := range obj.Pairs {
		result[v.Key.(NativeObject).Native()] = v.Value.(NativeObject).Native()
	}
	return result
}
