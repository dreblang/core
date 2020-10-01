package evaluator

import (
	"github.com/dreblang/core/object"
)

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func extendedFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendedFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		if result := fn.Fn(args...); result != nil {
			return result
		} else {
			return Null
		}

	case *object.MemberFn:
		if result := fn.Fn(fn.Obj, args...); result != nil {
			return result
		} else {
			return Null
		}

	default:
		return newError("%s: %s", notFunctionError, fn.Type())
	}
}
