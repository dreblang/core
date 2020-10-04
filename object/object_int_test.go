package object

import "testing"

func TestIntegerOperations(t *testing.T) {
	tests := []struct {
		num1     *Integer
		num2     Object
		op       string
		expected interface{}
	}{
		{&Integer{Value: 1}, &Integer{Value: 2}, "+", 3},
		{&Integer{Value: 1}, &Integer{Value: 2}, "-", -1},
		{&Integer{Value: 1}, &Integer{Value: 2}, "*", 2},
		{&Integer{Value: 1}, &Integer{Value: 2}, "/", 0},
		{&Integer{Value: 1}, &Integer{Value: 2}, "%", 1},

		{&Integer{Value: 1}, &Float{Value: 2}, "+", 3.0},
		{&Integer{Value: 1}, &Float{Value: 2}, "-", -1.0},
		{&Integer{Value: 1}, &Float{Value: 2}, "*", 2.0},
		{&Integer{Value: 1}, &Float{Value: 2}, "/", 0.5},

		{&Integer{Value: 1}, &Integer{Value: 2}, "<", true},
		{&Integer{Value: 1}, &Integer{Value: 2}, "<=", true},
		{&Integer{Value: 1}, &Integer{Value: 2}, ">", false},
		{&Integer{Value: 1}, &Integer{Value: 2}, ">=", false},
		{&Integer{Value: 1}, &Integer{Value: 2}, "==", false},
		{&Integer{Value: 1}, &Integer{Value: 2}, "!=", true},

		{&Integer{Value: 1}, &Float{Value: 2}, "<", true},
		{&Integer{Value: 1}, &Float{Value: 2}, "<=", true},
		{&Integer{Value: 1}, &Float{Value: 2}, ">", false},
		{&Integer{Value: 1}, &Float{Value: 2}, ">=", false},
		{&Integer{Value: 1}, &Float{Value: 1}, "==", true},
		{&Integer{Value: 1}, &Float{Value: 2}, "!=", true},

		{&Integer{Value: 1}, True, "==", false},
		{&Integer{Value: 1}, False, "!=", true},
	}

	for _, tt := range tests {
		result := tt.num1.InfixOperation(tt.op, tt.num2)
		testExpectForInt(t, result, tt.expected)
	}
}

func testExpectForInt(t *testing.T, obj Object, expected interface{}) bool {
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

func TestIntegerHashKey(t *testing.T) {
	hello1 := &Integer{Value: 1}
	hello2 := &Integer{Value: 1}
	diff1 := &Integer{Value: 2}
	diff2 := &Integer{Value: 2}

	if hello1.HashKey() != hello2.HashKey() ||
		diff1.HashKey() != diff2.HashKey() ||
		hello1.HashKey() == diff1.HashKey() {
		t.Errorf("ints with same content have different hash keys")
	}
}

func TestHelperFuncs(t *testing.T) {
	i := &Integer{Value: 1}

	if i.Inspect() != "1" {
		t.Errorf("Unexpected inspect result")
	}

	if i.String() != "1" {
		t.Errorf("Unexpected inspect result")
	}

	if i.GetMember("something").Type() != ErrorObj {
		t.Errorf("Expected error in get member")
	}

	if i.InfixOperation("+", True).Type() != ErrorObj {
		t.Errorf("Expected error in operation")
	}
}
