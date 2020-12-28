package math

import (
	"github.com/dreblang/core/object"

	"math"
)

func _floatVal(obj object.Object) (float64, bool) {
	switch v := obj.(type) {
	case *object.Float:
		return v.Value, true
	case *object.Integer:
		return float64(v.Value), true
	}
	return 0, false
}

var mathScope = &object.Scope{
	Name: "math",
	Exports: map[string]object.Object{
		"PI": &object.Float{Value: math.Pi},
		"E":  &object.Float{Value: math.E},

		"ceil":  &object.Builtin{Fn: ceil},
		"floor": &object.Builtin{Fn: floor},

		"fact": &object.Builtin{Fn: fact},

		"sin":  &object.Builtin{Fn: sin},
		"cos":  &object.Builtin{Fn: cos},
		"tan":  &object.Builtin{Fn: tan},
		"asin": &object.Builtin{Fn: asin},
		"acos": &object.Builtin{Fn: acos},
		"atan": &object.Builtin{Fn: atan},

		"sinh":  &object.Builtin{Fn: sinh},
		"cosh":  &object.Builtin{Fn: cosh},
		"tanh":  &object.Builtin{Fn: tanh},
		"asinh": &object.Builtin{Fn: asinh},
		"acosh": &object.Builtin{Fn: acosh},
		"atanh": &object.Builtin{Fn: atanh},

		"log":   &object.Builtin{Fn: log},
		"loge":  &object.Builtin{Fn: loge},
		"log2":  &object.Builtin{Fn: log2},
		"log10": &object.Builtin{Fn: log10},

		"deg2rad": &object.Builtin{Fn: deg2rad},
		"rad2deg": &object.Builtin{Fn: rad2deg},
	},
}

// Load the scope
func Load() *object.Scope {
	return mathScope
}

func ceil(args ...object.Object) object.Object {
	if x, ok := args[0].(*object.Float); ok {
		return &object.Float{Value: math.Ceil(x.Value)}
	} else if x, ok := args[0].(*object.Integer); ok {
		return x
	}
	return object.NullObject
}

func floor(args ...object.Object) object.Object {
	if x, ok := args[0].(*object.Float); ok {
		return &object.Float{Value: math.Floor(x.Value)}
	} else if x, ok := args[0].(*object.Integer); ok {
		return x
	}

	return object.NullObject
}

var _facts [41]int64

func calcFact(n int64) (res int64) {
	if _facts[n] != 0 {
		res = _facts[n]
		return res
	}

	if n > 0 {
		res = n * calcFact(n-1)
		_facts[n] = res
		return res
	}

	return 1
}

func fact(args ...object.Object) object.Object {
	if x, ok := args[0].(*object.Integer); ok {
		return &object.Integer{Value: calcFact(x.Value)}
	}
	return object.NullObject
}

func exp(args ...object.Object) object.Object {
	if x, ok := args[0].(*object.Float); ok {
		return &object.Float{Value: math.Exp(x.Value)}
	} else if x, ok := args[0].(*object.Integer); ok {
		return &object.Float{Value: math.Exp(float64(x.Value))}
	}
	return object.NullObject
}

func sin(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: math.Sin(x)}
		}
	}
	return object.NullObject
}
func cos(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: math.Cos(x)}
		}
	}
	return object.NullObject
}
func tan(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: math.Tan(x)}
		}
	}
	return object.NullObject
}
func asin(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: math.Asin(x)}
		}
	}
	return object.NullObject
}
func acos(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: math.Acos(x)}
		}
	}
	return object.NullObject
}
func atan(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: math.Atan(x)}
		}
	}
	return object.NullObject
}

func sinh(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: math.Sinh(x)}
		}
	}
	return object.NullObject
}
func cosh(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: math.Cosh(x)}
		}
	}
	return object.NullObject
}
func tanh(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: math.Tanh(x)}
		}
	}
	return object.NullObject
}
func asinh(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: math.Asinh(x)}
		}
	}
	return object.NullObject
}
func acosh(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: math.Acosh(x)}
		}
	}
	return object.NullObject
}
func atanh(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: math.Atanh(x)}
		}
	}
	return object.NullObject
}

func log(args ...object.Object) object.Object {
	if len(args) == 2 {
		if x, ok := _floatVal(args[0]); ok {
			if y, ok := _floatVal(args[1]); ok {
				return &object.Float{Value: math.Log(x) / math.Log(y)}
			}
		}
	}
	return object.NullObject
}
func loge(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: math.Log(x)}
		}
	}
	return object.NullObject
}
func log2(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: math.Log2(x)}
		}
	}
	return object.NullObject
}
func log10(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: math.Log10(x)}
		}
	}
	return object.NullObject
}
func pow(args ...object.Object) object.Object {
	if len(args) == 2 {
		if x, ok := _floatVal(args[0]); ok {
			if y, ok := _floatVal(args[1]); ok {
				return &object.Float{Value: math.Pow(x, y)}
			}
		}
	}
	return object.NullObject
}

func deg2rad(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: x * math.Pi / 180}
		}
	}
	return object.NullObject
}
func rad2deg(args ...object.Object) object.Object {
	if len(args) == 1 {
		if x, ok := _floatVal(args[0]); ok {
			return &object.Float{Value: x / math.Pi * 180}
		}
	}
	return object.NullObject
}
