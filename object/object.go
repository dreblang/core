package object

type ObjectType string

const (
	IntegerObj          = "Integer"
	FloatObj            = "Float"
	BooleanObj          = "Boolean"
	NullObj             = "Null"
	ReturnValueObj      = "ReturnValue"
	ErrorObj            = "Error"
	FunctionObj         = "Function"
	BytesObj            = "Bytes"
	StringObj           = "String"
	BuiltinObj          = "Builtin"
	ArrayObj            = "Array"
	HashObj             = "Hash"
	CompiledFunctionObj = "CompiledFunction"
	ClosureObj          = "Closure"
	ScopeObj            = "Scope"
	ClassObj            = "Class"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	String() string

	GetMember(name string) Object
	SetMember(name string, value Object) Object

	Equals(other Object) bool
}

type NativeObject interface {
	Object
	Native() interface{}
}

type InfixOperatorObject interface {
	Object
	InfixOperation(operator string, other Object) Object
}
