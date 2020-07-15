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
	StringObj           = "String"
	BuiltinObj          = "Builtin"
	ArrayObj            = "Array"
	HashObj             = "Hash"
	CompiledFunctionObj = "CompiledFunction"
	ClosureObj          = "Closure"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}
