package evaluator

import (
	"log"

	"github.com/dreblang/core/ast"
	"github.com/dreblang/core/object"
)

const (
	unknownOperatorError    = "unknown eval operator"
	typeMissMatchError      = "type mismatch"
	identifierNotFoundError = "identifier not found"
	notFunctionError        = "not a function"
)

var (
	Null  = object.NullObject
	True  = object.True
	False = object.False
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	if node == nil {
		return nil
	}

	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.LoopExpression:
		return evalLoopExpression(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		if node.Operator == "." {
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			return evalMemberOperation(left, node.Right)

		} else if node.Operator == "=" {
			val := Eval(node.Right, env)
			if isError(val) {
				return val
			}
			switch leftNode := node.Left.(type) {
			case *ast.Identifier:
				env.Set(leftNode.String(), val)
				return val

			case *ast.IndexExpression:
				return evalIndexSetExpression(
					Eval(leftNode.Left, env),
					Eval(leftNode.Index, env),
					Eval(leftNode.IndexUpper, env),
					Eval(leftNode.IndexSkip, env),
					leftNode.HasUpper,
					leftNode.HasSkip,
					val,
				)

			default:
				return newError("Error executing assignment.")
			}

		} else {
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return evalInfixExpression(node.Operator, left, right)
		}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		var index, indexUpper, indexSkip object.Object
		if node.Index != nil {
			index = Eval(node.Index, env)
		}
		if node.IndexUpper != nil {
			indexUpper = Eval(node.IndexUpper, env)
		}
		if node.IndexSkip != nil {
			indexSkip = Eval(node.IndexSkip, env)
		}
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index, indexUpper, indexSkip, node.HasUpper, node.HasSkip)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return object.NativeBoolToBooleanObject(node.Value)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	default:
		log.Println("Unknown expression!")
	}

	return nil
}
