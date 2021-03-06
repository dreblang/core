package vm

import (
	"errors"
	"fmt"

	"github.com/dreblang/core/code"
	"github.com/dreblang/core/compiler"
	"github.com/dreblang/core/object"
)

const StackSize = 4096
const GlobalSize = 65536
const MaxFrames = 2048

var True = object.True
var False = object.False
var Null = object.NullValue

var stackOverflowErr = errors.New("stack overflow")
var frameOverflowErr = errors.New("frame overflow")

type VM struct {
	constants   []object.Object
	stack       []object.Object
	sp          int // Always points to the next value. Top of stack is stack[sp-1]
	globals     []object.Object
	frames      []*Frame
	framesIndex int

	curFrame *Frame
}

var currentVM *VM

func GetCurrentVM() *VM {
	return currentVM
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainClosure := &object.Closure{Fn: mainFn, Exports: map[string]object.Object{}}
	mainFrame := NewFrame(mainClosure, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	return &VM{
		constants:   bytecode.Constants,
		stack:       make([]object.Object, StackSize),
		sp:          0,
		globals:     make([]object.Object, GlobalSize),
		frames:      frames,
		framesIndex: 1,
		curFrame:    mainFrame,
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = s
	return vm
}

func (vm *VM) Run() error {
	var ip int
	var ins code.Instructions
	var op code.Opcode

	currentVM = vm

	for vm.curFrame.ip < len(vm.curFrame.instructions)-1 {
		var err error
		vm.curFrame.ip++

		ip = vm.curFrame.ip
		ins = vm.curFrame.instructions
		op = code.Opcode(ins[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(ins[ip+1:])
			vm.curFrame.ip += 2

			err = vm.push(vm.constants[constIndex])
		case code.OpPop:
			vm.pop()

		case code.OpAdd:
			err = vm.executeBinaryOperation("+")
		case code.OpSub:
			err = vm.executeBinaryOperation("-")
		case code.OpMul:
			err = vm.executeBinaryOperation("*")
		case code.OpDiv:
			err = vm.executeBinaryOperation("/")
		case code.OpMod:
			err = vm.executeBinaryOperation("%")
		case code.OpMember:
			err = vm.executeMemberOperation(op)
		case code.OpMemberSet:
			err = vm.executeMemberSetOperation(op)
		case code.OpTrue:
			err = vm.push(True)
		case code.OpFalse:
			err = vm.push(False)
		case code.OpEqual:
			err = vm.executeComparison("==")
		case code.OpNotEqual:
			err = vm.executeComparison("!=")
		case code.OpGreaterThan:
			err = vm.executeComparison(">")
		case code.OpGreaterOrEqual:
			err = vm.executeComparison(">=")
		case code.OpBang:
			err = vm.executeBangOperator()
		case code.OpMinus:
			err = vm.executeMinusOperator()
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.curFrame.ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				vm.curFrame.ip = pos - 1
			}
		case code.OpJump:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.curFrame.ip = pos - 1
		case code.OpNull:
			err = vm.push(Null)
		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.curFrame.ip += 2

			val := vm.pop()
			vm.globals[globalIndex] = val
			err = vm.push(val)
		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.curFrame.ip += 2

			err = vm.push(vm.globals[globalIndex])
		case code.OpArray:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.curFrame.ip += 2

			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements

			err = vm.push(array)
		case code.OpHash:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.curFrame.ip += 2

			hash, err := vm.buildHash(vm.sp-numElements, vm.sp)
			if err != nil {
				return err
			}
			vm.sp = vm.sp - numElements

			err = vm.push(hash)
		case code.OpIndex:
			hasSkip := vm.pop()
			indexSkip := vm.pop()
			hasUpper := vm.pop()
			indexUpper := vm.pop()
			index := vm.pop()
			left := vm.pop()
			err = vm.executeIndexExpression(left, index, indexUpper, indexSkip, hasUpper, hasSkip)

		case code.OpIndexSet:
			hasSkip := vm.pop()
			indexSkip := vm.pop()
			hasUpper := vm.pop()
			indexUpper := vm.pop()
			index := vm.pop()
			left := vm.pop()
			right := vm.pop()
			err = vm.executeIndexSetExpression(left, index, indexUpper, indexSkip, hasUpper, hasSkip, right)

		case code.OpCall:
			numArgs := code.ReadUint8(ins[ip+1:])
			vm.curFrame.ip++

			err = vm.executeCall(int(numArgs))
		case code.OpReturnValue:
			returnValue := vm.pop()

			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err = vm.push(returnValue)
		case code.OpReturn:
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err = vm.push(Null)
		case code.OpSetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.curFrame.ip++

			frame := vm.curFrame
			val := vm.pop()
			vm.stack[frame.basePointer+int(localIndex)] = val
			err = vm.push(val)
		case code.OpGetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.curFrame.ip++

			frame := vm.curFrame

			err = vm.push(vm.stack[frame.basePointer+int(localIndex)])
		case code.OpGetBuiltin:
			builtinIndex := code.ReadUint8(ins[ip+1:])
			vm.curFrame.ip++

			definition := object.Builtins[builtinIndex]

			err = vm.push(definition.Builtin)
		case code.OpClosure:
			constIndex := code.ReadUint16(ins[ip+1:])
			numFree := code.ReadUint8(ins[ip+3:])
			vm.curFrame.ip += 3

			err = vm.pushClosure(int(constIndex), int(numFree))

		case code.OpSetFree:
			freeIndex := code.ReadUint8(ins[ip+1:])
			vm.curFrame.ip++

			currentClosure := vm.curFrame.cl
			val := vm.pop()
			currentClosure.Free[freeIndex] = val
			err = vm.push(val)

		case code.OpGetFree:
			freeIndex := code.ReadUint8(ins[ip+1:])
			vm.curFrame.ip++

			currentClosure := vm.curFrame.cl
			err = vm.push(currentClosure.Free[freeIndex])

		case code.OpScope:
			scopeObj := &object.Scope{
				Exports: vm.curFrame.cl.Exports,
			}
			vm.push(scopeObj)

		case code.OpScopeResolve:
			name := vm.pop()
			scope := vm.pop()
			callable := scope.GetMember(name.String())
			vm.push(callable)

		case code.OpExport:
			val := vm.pop()
			name := vm.pop()
			vm.curFrame.cl.Exports[name.(*object.String).Value] = val
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return stackOverflowErr
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	if vm.sp == 0 {
		return object.NullObject
	}

	o := vm.stack[vm.sp-1]
	vm.sp--
	return o
}

func (vm *VM) executeBinaryOperation(op string) error {
	right := vm.pop()
	left := vm.pop()

	if leftOp, ok := left.(object.InfixOperatorObject); ok {
		result := leftOp.InfixOperation(op, right)
		if result.Type() != object.ErrorObj {
			vm.push(result)
			return nil
		} else {
			return fmt.Errorf(result.(*object.Error).Message)
		}
	} else {
		return fmt.Errorf("%s: %s %s %s", "unknown eval operator", left.Type(), op, right.Type())
	}
}

func (vm *VM) executeMemberOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	result := left.GetMember(right.String())
	vm.push(result)
	return nil
}

func (vm *VM) executeMemberSetOperation(op code.Opcode) error {
	left := vm.pop()
	member := vm.pop()
	right := vm.pop()

	result := left.SetMember(member.String(), right)
	if result.Type() != object.ErrorObj {
		vm.push(result)
		return nil
	}

	return fmt.Errorf(result.(*object.Error).Message)
}

func (vm *VM) executeComparison(op string) error {
	right := vm.pop()
	left := vm.pop()

	if leftOp, ok := left.(object.InfixOperatorObject); ok {
		result := leftOp.InfixOperation(op, right)
		if result.Type() != object.ErrorObj {
			vm.push(result)
			return nil
		} else {
			return fmt.Errorf(result.(*object.Error).Message)
		}
	} else {
		return fmt.Errorf("%s: %s %s %s", "unknownOperatorError", left.Type(), op, right.Type())
	}
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()

	switch val := operand.(type) {
	case *object.Integer:
		return vm.push(&object.Integer{Value: -val.Value})
	case *object.Float:
		return vm.push(&object.Float{Value: -val.Value})
	}

	return fmt.Errorf("unsupported type for negation: %s", operand.Type())
}

func (vm *VM) executeIndexExpression(left, index, indexUpper, indexSkip, hasUpper, hasSkip object.Object) error {
	switch {
	case left.Type() == object.ArrayObj:
		return vm.executeArrayIndex(left, index, indexUpper, indexSkip, hasUpper, hasSkip)

	case left.Type() == object.HashObj:
		return vm.executeHashIndex(left, index)
	default:
		return fmt.Errorf("index operator not supported: %s", left.Type())
	}
}

func (vm *VM) executeIndexSetExpression(left, index, indexUpper, indexSkip, hasUpper, hasSkip, right object.Object) error {
	switch {
	case left.Type() == object.ArrayObj:
		return vm.executeArrayIndexSet(left, index, indexUpper, indexSkip, hasUpper, hasSkip, right)

	case left.Type() == object.HashObj:
		return vm.executeHashIndexSet(left, index, right)
	default:
		return fmt.Errorf("index set operator not supported: %s", left.Type())
	}
}

func (vm *VM) executeArrayIndex(array, index, indexUpper, indexSkip, hasUpper, hasSkip object.Object) error {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements))

	var idxUpper int64 = max

	if idx < 0 {
		idx = max + idx
	}

	if !isTruthy(hasUpper) {
		if idx < 0 || idx >= max {
			return vm.push(Null)
		}

		return vm.push(arrayObject.Elements[idx])
	}

	if indexUpper != object.NullObject {
		idxUpper = indexUpper.(*object.Integer).Value
	}
	if idxUpper < 0 {
		idxUpper += max
	}

	var inc int64 = 1
	if isTruthy(hasSkip) {
		inc = indexSkip.(*object.Integer).Value
	}

	elements := make([]object.Object, 0)

	for i := idx; i < idxUpper; i += inc {
		elements = append(elements, arrayObject.Elements[i])
	}
	return vm.push(&object.Array{
		Elements: elements,
	})
}

func (vm *VM) executeArrayIndexSet(array, index, indexUpper, indexSkip, hasUpper, hasSkip, right object.Object) error {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements))

	// var idxUpper int64 = max

	if idx < 0 {
		idx = max + idx
	}

	if !isTruthy(hasUpper) {
		if idx >= max {
			return fmt.Errorf("index out of bounds")
		}
		arrayObject.Elements[idx] = right
		return nil
	}
	return fmt.Errorf("Cannot assign range!")
	// TODO: Implement setting of array
	// if indexUpper != object.NullObject {
	// 	idxUpper = indexUpper.(*object.Integer).Value
	// }
	// if idxUpper < 0 {
	// 	idxUpper += max
	// }

	// var inc int64 = 1
	// if isTruthy(hasSkip) {
	// 	inc = indexSkip.(*object.Integer).Value
	// }

	// elements := make([]object.Object, 0)

	// for i := idx; i < idxUpper; i += inc {
	// 	elements = append(elements, arrayObject.Elements[i])
	// }
	// return vm.push(&object.Array{
	// 	Elements: elements,
	// })
}

func (vm *VM) executeHashIndex(hash, index object.Object) error {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return fmt.Errorf("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return vm.push(Null)
	}

	return vm.push(pair.Value)
}

func (vm *VM) executeHashIndexSet(hash, index, right object.Object) error {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)
	if !ok {
		return fmt.Errorf("unusable as hash key: %s", index.Type())
	}

	hashObject.Pairs[key.HashKey()] = object.HashPair{
		Key:   index,
		Value: right,
	}
	return nil
}

func (vm *VM) buildArray(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return &object.Array{Elements: elements}
}

func (vm *VM) buildHash(startIndex, endIndex int) (object.Object, error) {
	hashedPairs := make(map[object.HashKey]object.HashPair)

	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]

		pair := object.HashPair{Key: key, Value: value}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("unusable as hash key: %s", key.Type())
		}

		hashedPairs[hashKey.HashKey()] = pair
	}

	return &object.Hash{Pairs: hashedPairs}, nil
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushFrame(f *Frame) error {
	if vm.framesIndex >= MaxFrames {
		return frameOverflowErr
	}
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++

	vm.curFrame = vm.currentFrame()
	return nil
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	vm.curFrame = vm.currentFrame()

	return vm.frames[vm.framesIndex]
}

func (vm *VM) executeCall(numArgs int) error {
	callee := vm.stack[vm.sp-1-numArgs]
	switch callee := callee.(type) {
	case *object.Closure:
		return vm.callClosure(callee, numArgs)
	case *object.Builtin:
		return vm.callBuiltin(callee, numArgs)
	case *object.MemberFn:
		return vm.callMember(callee, numArgs)
	default:
		return fmt.Errorf("calling non-function and non-built-in")
	}
}

func (vm *VM) callClosure(cl *object.Closure, numArgs int) error {
	if numArgs != cl.Fn.NumParameters {
		return fmt.Errorf("wrong number of arguments: want=%d, got=%d",
			cl.Fn.NumParameters, numArgs)
	}

	frame := NewFrame(cl, vm.sp-numArgs)
	err := vm.pushFrame(frame)
	if err != nil {
		return err
	}

	vm.sp = frame.basePointer + cl.Fn.NumLocals

	return nil
}

func (vm *VM) callBuiltin(builtin *object.Builtin, numArgs int) error {
	args := vm.stack[vm.sp-numArgs : vm.sp]

	result := builtin.Fn(args...)
	vm.sp = vm.sp - numArgs - 1

	if result != nil {
		vm.push(result)
	} else {
		vm.push(Null)
	}

	return nil
}

func (vm *VM) callMember(memberfn *object.MemberFn, numArgs int) error {
	args := vm.stack[vm.sp-numArgs : vm.sp]

	result := memberfn.Fn(memberfn.Obj, args...)
	vm.sp = vm.sp - numArgs - 1

	if result != nil {
		vm.push(result)
	} else {
		vm.push(Null)
	}

	return nil
}

func (vm *VM) pushClosure(constIndex, numFree int) error {
	constant := vm.constants[constIndex]
	function := constant.(*object.CompiledFunction)

	free := make([]object.Object, numFree)
	for i := 0; i < numFree; i++ {
		free[i] = vm.stack[vm.sp-numFree+i]
	}
	vm.sp = vm.sp - numFree

	closure := &object.Closure{Fn: function, Free: free, Exports: map[string]object.Object{}}
	return vm.push(closure)
}

func (vm *VM) ExecClosure(closure *object.Closure, args ...object.Object) object.Object {
	nvm := NewWithGlobalsStore(
		&compiler.Bytecode{
			Instructions: code.Instructions{},
			Constants:    vm.constants,
		},
		vm.globals,
	)
	nvm.push(Null)
	for _, arg := range args {
		nvm.push(arg)
	}
	err := nvm.callClosure(closure, len(args))

	if err != nil {
		return object.NewError("Error calling closure: %s", err)
	}

	nvm.Run()
	if nvm.sp >= 0 {
		return nvm.stack[nvm.sp-1]
	}
	return Null
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}
