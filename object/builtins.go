package object

import (
	"fmt"
	"strconv"
)

const (
	BuiltinFuncNameLen   = "len"
	BuiltinFuncNamePrint = "print"

	// Type conversions
	BuiltinFuncNameInt    = "int"
	BuiltinFuncNameFloat  = "float"
	BuiltinFuncNameString = "string"
)

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		BuiltinFuncNameLen,
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to %q not supported, got %s",
					BuiltinFuncNameLen, args[0].Type())
			}
		},
		},
	},
	{
		BuiltinFuncNamePrint,
		&Builtin{Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}
			fmt.Println()

			return nil
		},
		},
	},
	{
		BuiltinFuncNameInt,
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			switch arg := args[0].(type) {
			case *Integer:
				return arg
			case *Float:
				return &Integer{Value: int64(arg.Value)}
			case *String:
				val, err := strconv.ParseInt(arg.Value, 10, 64)
				if err != nil {
					return newError("Conversion to int failed!")
				}
				return &Integer{Value: val}
			default:
				return newError("argument to %q not supported, got %s",
					BuiltinFuncNameLen, args[0].Type())
			}
		},
		},
	},
	{
		BuiltinFuncNameFloat,
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			switch arg := args[0].(type) {
			case *Integer:
				return &Float{Value: float64(arg.Value)}
			case *Float:
				return arg
			case *String:
				val, err := strconv.ParseFloat(arg.Value, 64)
				if err != nil {
					return newError("Conversion to float failed!")
				}
				return &Float{Value: val}
			default:
				return newError("argument to %q not supported, got %s",
					BuiltinFuncNameLen, args[0].Type())
			}
		},
		},
	},
	{
		BuiltinFuncNameString,
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			return &String{Value: args[0].String()}
		},
		},
	},
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func NewError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
