package object

import "testing"

func TestFloatOperations(t *testing.T) {
	tests := []struct {
		num1     *Float
		num2     Object
		op       string
		expected interface{}
	}{
		{&Float{Value: 1}, &Integer{Value: 2}, "+", 3.0},
		{&Float{Value: 1}, &Integer{Value: 2}, "-", -1.0},
		{&Float{Value: 1}, &Integer{Value: 2}, "*", 2.0},
		{&Float{Value: 1}, &Integer{Value: 2}, "/", 0.5},

		{&Float{Value: 1}, &Float{Value: 2}, "+", 3.0},
		{&Float{Value: 1}, &Float{Value: 2}, "-", -1.0},
		{&Float{Value: 1}, &Float{Value: 2}, "*", 2.0},
		{&Float{Value: 1}, &Float{Value: 2}, "/", 0.5},

		{&Float{Value: 1}, &Integer{Value: 2}, "<", true},
		{&Float{Value: 1}, &Integer{Value: 2}, "<=", true},
		{&Float{Value: 1}, &Integer{Value: 1}, "<=", true},
		{&Float{Value: 1}, &Integer{Value: 2}, ">", false},
		{&Float{Value: 1}, &Integer{Value: 2}, ">=", false},
		{&Float{Value: 1}, &Integer{Value: 1}, ">=", true},
		{&Float{Value: 1}, &Integer{Value: 2}, "==", false},
		{&Float{Value: 1}, &Integer{Value: 2}, "!=", true},

		{&Float{Value: 1}, &Float{Value: 2}, "<", true},
		{&Float{Value: 1}, &Float{Value: 2}, "<=", true},
		{&Float{Value: 1}, &Float{Value: 2}, ">", false},
		{&Float{Value: 1}, &Float{Value: 2}, ">=", false},
		{&Float{Value: 1}, &Float{Value: 1}, "==", true},
		{&Float{Value: 1}, &Float{Value: 2}, "!=", true},

		{&Float{Value: 1}, True, "==", false},
		{&Float{Value: 1}, False, "!=", true},
	}

	for _, tt := range tests {
		result := tt.num1.InfixOperation(tt.op, tt.num2)
		testExpectForFloat(t, result, tt.expected)
	}
}

func testExpectForFloat(t *testing.T, obj Object, expected interface{}) bool {
	if intNum, ok := expected.(int); ok {
		result, ok := obj.(*Integer)
		if !ok {
			t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
			return false
		}

		if result.Value != int64(intNum) {
			t.Errorf("object has wrong value. got=%d, want=%d", result.Value, intNum)
			return false
		}

		return true

	} else if floatNum, ok := expected.(float64); ok {
		result, ok := obj.(*Float)
		if !ok {
			t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
			return false
		}

		if result.Value != floatNum {
			t.Errorf("object has wrong value. got=%v, want=%v", result.Value, floatNum)
			return false
		}

		return true

	} else if boolVal, ok := expected.(bool); ok {
		result, ok := obj.(*Boolean)
		if !ok {
			t.Errorf("object is not Bool. got=%T (%+v)", obj, obj)
			return false
		}

		if result.Value != boolVal {
			t.Errorf("object has wrong value. got=%v, want=%v", result.Value, floatNum)
			return false
		}

		return true
	}

	t.Errorf("Unknown expectation! (%+v)", expected)
	return false
}

func TestFloatHashKey(t *testing.T) {
	hello1 := &Float{Value: 1.1}
	hello2 := &Float{Value: 1.1}
	diff1 := &Float{Value: 2.2}
	diff2 := &Float{Value: 2.2}

	if hello1.HashKey() != hello2.HashKey() ||
		diff1.HashKey() != diff2.HashKey() ||
		hello1.HashKey() == diff1.HashKey() {
		t.Errorf("floats with same content have different hash keys")
	}
}

func TestFloatHelperFuncs(t *testing.T) {
	i := &Float{Value: 1}

	if i.Inspect() != "1.000000" {
		t.Errorf("Unexpected inspect result %s", i.Inspect())
	}

	if i.String() != "1" {
		t.Errorf("Unexpected inspect result %s", i.String())
	}

	if i.GetMember("something").Type() != ErrorObj {
		t.Errorf("Expected error in get member")
	}

	if i.InfixOperation("+", True).Type() != ErrorObj {
		t.Errorf("Expected error in operation")
	}
}
