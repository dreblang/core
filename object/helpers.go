package object

var (
	NullValue = NullObject
	True      = &Boolean{Value: true}
	False     = &Boolean{Value: false}
)

const (
	unknownOperatorError = "unknown eval operator"
	typeMissMatchError   = "type mismatch"
)

func NativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return True
	} else {
		return False
	}
}
