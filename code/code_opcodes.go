package code

type Opcode byte

const (
	OpConstant Opcode = iota

	OpPop

	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod

	OpTrue
	OpFalse

	OpEqual
	OpNotEqual
	OpGreaterThan
	OpGreaterOrEqual

	OpMinus
	OpBang

	OpJumpNotTruthy
	OpJump

	OpNull

	OpGetGlobal
	OpSetGlobal
	OpGetLocal
	OpSetLocal
	OpGetFree
	OpSetFree

	OpArray
	OpHash

	OpIndex
	OpIndexSet

	OpCall
	OpReturnValue
	OpReturn

	OpGetBuiltin
	OpClosure
	OpMember

	OpExport
	OpScope
	OpScopeResolve
)
