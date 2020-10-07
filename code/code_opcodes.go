package code

type Opcode int

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
	OpArray
	OpHash
	OpIndex
	OpIndexSet
	OpCall
	OpReturnValue
	OpReturn
	OpGetLocal
	OpSetLocal
	OpGetBuiltin
	OpClosure
	OpGetFree
	OpSetFree
	OpMember
	OpExport
	OpScope
	OpScopeResolve
)
