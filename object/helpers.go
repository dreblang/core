package object

var (
	NullValue = &Null{}
	True      = &Boolean{Value: true}
	False     = &Boolean{Value: false}
)

const (
	unknownOperatorError = "unknown eval operator"
)

func NativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return True
	} else {
		return False
	}
}
