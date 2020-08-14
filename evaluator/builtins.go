package evaluator

import (
	"github.com/dreblang/core/object"
)

var builtins = map[string]*object.Builtin{
	object.BuiltinFuncNameLen:   object.GetBuiltinByName(object.BuiltinFuncNameLen),
	object.BuiltinFuncNamePuts:  object.GetBuiltinByName(object.BuiltinFuncNamePuts),
	object.BuiltinFuncNameFirst: object.GetBuiltinByName(object.BuiltinFuncNameFirst),
	object.BuiltinFuncNameLast:  object.GetBuiltinByName(object.BuiltinFuncNameLast),
	object.BuiltinFuncNameRest:  object.GetBuiltinByName(object.BuiltinFuncNameRest),
	object.BuiltinFuncNamePush:  object.GetBuiltinByName(object.BuiltinFuncNamePush),

	object.BuiltinFuncNameInt:    object.GetBuiltinByName(object.BuiltinFuncNameInt),
	object.BuiltinFuncNameFloat:  object.GetBuiltinByName(object.BuiltinFuncNameFloat),
	object.BuiltinFuncNameString: object.GetBuiltinByName(object.BuiltinFuncNameString),
}
