package object

import (
	"fmt"
	"strconv"
)

const (
	BuiltinFuncNameLen   = "len"
	BuiltinFuncNameFirst = "first"
	BuiltinFuncNameLast  = "last"
	BuiltinFuncNameRest  = "rest"
	BuiltinFuncNamePush  = "push"
	BuiltinFuncNamePuts  = "puts"

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
		BuiltinFuncNamePuts,
		&Builtin{Fn: func(args ...Object) Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return nil
		},
		},
	},
	{
		BuiltinFuncNameFirst,
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != ArrayObj {
				return newError("argument to %q must be %s, got %s",
					BuiltinFuncNameFirst, ArrayObj, args[0].Type())
			}

			arr := args[0].(*Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return nil
		},
		},
	},
	{
		BuiltinFuncNameLast,
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != ArrayObj {
				return newError("argument to %q must be %s, got %s",
					BuiltinFuncNameLast, ArrayObj, args[0].Type())
			}

			arr := args[0].(*Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return nil
		},
		},
	},
	{
		BuiltinFuncNameRest,
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != ArrayObj {
				return newError("argument to %q must be %s, got %s",
					BuiltinFuncNameRest, ArrayObj, args[0].Type())
			}

			arr := args[0].(*Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &Array{Elements: newElements}
			}

			return nil
		},
		},
	},
	{
		BuiltinFuncNamePush,
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2",
					len(args))
			}
			if args[0].Type() != ArrayObj {
				return newError("argument to %q must be %s, got %s",
					BuiltinFuncNamePush, ArrayObj, args[0].Type())
			}

			arr := args[0].(*Array)
			length := len(arr.Elements)

			newElements := make([]Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &Array{Elements: newElements}
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

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}
