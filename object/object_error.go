package object

import "errors"

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ErrorObj }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) String() string   { return e.Message }

func (obj *Error) GetMember(name string) Object {
	return newError("No member named [%s]", name)
}

func (obj *Error) SetMember(name string, value Object) Object {
	return newError("No member named [%s]", name)
}

func (obj *Error) Native() interface{} {
	return errors.New(obj.Message)
}

func (obj *Error) Equals(other Object) bool {
	if otherObj, ok := other.(*Error); ok {
		return obj.Message == otherObj.Message
	}
	return false
}
