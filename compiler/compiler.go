package compiler

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"plugin"
	"sort"
	"strings"

	"github.com/dreblang/core/ast"
	"github.com/dreblang/core/code"
	"github.com/dreblang/core/lexer"
	"github.com/dreblang/core/object"
	"github.com/dreblang/core/parser"
	"github.com/dreblang/core/token"
)

type Compiler struct {
	instructions        code.Instructions
	constants           []object.Object
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
	symbolTable         *SymbolTable
	scopes              []CompilationScope
	scopeIndex          int
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

func New() *Compiler {
	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	symbolTable := NewSymbolTable()

	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	return &Compiler{
		constants:   []object.Object{},
		symbolTable: symbolTable,
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
	}
}

func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = s
	compiler.constants = constants
	return compiler
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)

	case *ast.InfixExpression:
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.emit(code.OpGreaterThan)
			return nil
		} else if node.Operator == "<=" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.emit(code.OpGreaterOrEqual)
			return nil

		} else if node.Operator == "=" {
			var symbol Symbol
			switch leftNode := node.Left.(type) {
			case *ast.Identifier:
				symbol = c.symbolTable.Define(leftNode.String())
			}
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			switch leftNode := node.Left.(type) {
			case *ast.Identifier:
				c.saveSymbol(symbol)

			case *ast.IndexExpression:
				c.Compile(leftNode.Left)
				if leftNode.Index != nil {
					err = c.Compile(leftNode.Index)
					if err != nil {
						return err
					}
				} else {
					c.emit(code.OpConstant, c.addConstant(&object.Integer{}))
				}

				if leftNode.IndexUpper != nil {
					err = c.Compile(leftNode.IndexUpper)
					if err != nil {
						return err
					}
				} else {
					c.emit(code.OpNull)
				}
				c.emit(code.OpConstant, c.addConstant(object.NativeBoolToBooleanObject(leftNode.HasUpper)))

				if leftNode.IndexSkip != nil {
					err = c.Compile(leftNode.IndexSkip)
					if err != nil {
						return err
					}
				} else {
					c.emit(code.OpConstant, c.addConstant(&object.Integer{Value: 1}))
				}
				c.emit(code.OpConstant, c.addConstant(object.NativeBoolToBooleanObject(leftNode.HasSkip)))
				c.emit(code.OpIndexSet)

			case *ast.InfixExpression:
				if leftNode.Operator == "." {
					c.emit(code.OpConstant, c.addConstant(&object.String{Value: leftNode.Right.String()}))
					c.Compile(leftNode.Left)
					c.emit(code.OpMemberSet)
				}
			}

			return nil
		}

		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case "%":
			c.emit(code.OpMod)
		case ">":
			c.emit(code.OpGreaterThan)
		case ">=":
			c.emit(code.OpGreaterOrEqual)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		case ".":
			c.emit(code.OpMember)
		case "::":
			c.emit(code.OpScopeResolve)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// Emit an `OpJumpNotTruthy` with a bogus value
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}

		// Emit an `OpJump` with a bogus value
		jumpPos := c.emit(code.OpJump, 9999)

		afterConsequencePos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIs(code.OpPop) {
				c.removeLastPop()
			}
		}

		afterAlternative := len(c.currentInstructions())
		c.changeOperand(jumpPos, afterAlternative)

	case *ast.LoopExpression:
		blockStart := len(c.currentInstructions())
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// Emit an `OpJumpNotTruthy` with a bogus value
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		c.emit(code.OpJump, blockStart)

		afterConsequencePos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		if node.Index != nil {
			err = c.Compile(node.Index)
			if err != nil {
				return err
			}
		} else {
			c.emit(code.OpConstant, c.addConstant(&object.Integer{}))
		}

		if node.IndexUpper != nil {
			err = c.Compile(node.IndexUpper)
			if err != nil {
				return err
			}
		} else {
			c.emit(code.OpConstant, c.addConstant(object.NullObject))
		}
		c.emit(code.OpConstant, c.addConstant(object.NativeBoolToBooleanObject(node.HasUpper)))

		if node.IndexSkip != nil {
			err = c.Compile(node.IndexSkip)
			if err != nil {
				return err
			}
		} else {
			c.emit(code.OpConstant, c.addConstant(&object.Integer{Value: 1}))
		}
		c.emit(code.OpConstant, c.addConstant(object.NativeBoolToBooleanObject(node.HasSkip)))

		c.emit(code.OpIndex)

	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}

		for _, a := range node.Arguments {
			err := c.Compile(a)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpCall, len(node.Arguments))

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.LetStatement:
		symbol := c.symbolTable.Define(node.Name.Value)
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		c.saveSymbol(symbol)

	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}

		c.emit(code.OpReturnValue)

	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		c.loadSymbol(symbol)

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))

	case *ast.FloatLiteral:
		float := &object.Float{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(float))

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			err := c.Compile(el)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpArray, len(node.Elements))

	case *ast.HashLiteral:
		keys := []ast.Expression{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			err := c.Compile(k)
			if err != nil {
				err = c.Compile(
					&ast.StringLiteral{
						Token: token.Token{},
						Value: k.String(),
					},
				)
				if err != nil {
					return err
				}
			}
			err = c.Compile(node.Pairs[k])
			if err != nil {
				return err
			}
		}

		c.emit(code.OpHash, len(node.Pairs)*2)

	case *ast.FunctionLiteral:
		c.enterScope()

		for _, p := range node.Parameters {
			c.symbolTable.Define(p.Value)
		}

		err := c.Compile(node.Body)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastPopWithReturn()
		}
		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}

		freeSymbols := c.symbolTable.FreeSymbols
		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()

		for _, s := range freeSymbols {
			c.loadSymbol(s)
		}

		compiledFn := &object.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(node.Parameters),
		}

		fnIndex := c.addConstant(compiledFn)
		c.emit(code.OpClosure, fnIndex, len(freeSymbols))

	case *ast.ScopeDefinition:
		c.enterScope()

		err := c.Compile(node.Block)
		if err != nil {
			return err
		}
		c.emit(code.OpScope)
		c.emit(code.OpReturnValue)

		freeSymbols := c.symbolTable.FreeSymbols
		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()

		compiledFn := &object.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: 0,
		}
		fnIndex := c.addConstant(compiledFn)
		c.emit(code.OpClosure, fnIndex, len(freeSymbols))
		c.emit(code.OpCall, 0)

		symbol := c.symbolTable.Define(node.Name.Value)
		c.saveSymbol(symbol)

	case *ast.ExportStatement:
		nameConst := c.addConstant(&object.String{Value: node.Identifier.Value})
		c.emit(code.OpConstant, nameConst)
		symbol, ok := c.symbolTable.Resolve(node.Identifier.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Identifier.Value)
		}
		c.loadSymbol(symbol)
		c.emit(code.OpExport)

	case *ast.LoadStatement:
		c.loadModule(node.Identifier.Value)

	default:
		fmt.Println("Unknown ast type!", node)
	}

	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

func (c *Compiler) searchFile(fname string) *string {
	// Check in current working directory
	d, _ := os.Getwd()
	p := path.Join(d, fname)
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		return &p
	}

	if value, ok := os.LookupEnv("DREB_PATH"); ok {
		dirs := strings.Split(value, ":")
		for _, d := range dirs {
			p := path.Join(d, fname)
			if _, err := os.Stat(p); !os.IsNotExist(err) {
				return &p
			}
		}
	}

	return nil
}

func (c *Compiler) SearchPlugin(m string) *string {
	return c.searchFile(m + ".so")
}

func (c *Compiler) SearchSource(m string) *string {
	return c.searchFile(m + ".dreb")
}

func (c *Compiler) addConstant(obj object.Object) int {
	for idx, v := range c.constants {
		if obj.Equals(v) {
			return idx
		}
	}
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), ins...)

	c.scopes[c.scopeIndex].instructions = updatedInstructions

	return posNewInstruction
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.scopes[c.scopeIndex].lastInstruction.Opcode == code.OpPop
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	new := old[:last.Position]

	c.scopes[c.scopeIndex].instructions = new
	c.scopes[c.scopeIndex].lastInstruction = previous
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := c.currentInstructions()

	for i := 0; i < len(newInstruction); i++ {
		ins[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--
	c.symbolTable = c.symbolTable.Outer

	return instructions
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, code.Make(code.OpReturnValue))

	c.scopes[c.scopeIndex].lastInstruction.Opcode = code.OpReturnValue
}

func (c *Compiler) loadModule(m string) {
	var scope *object.Scope
	if loader, ok := coreModules[m]; ok {
		scope = loader()
	}

	if pluginFile := c.SearchPlugin(m); scope == nil && pluginFile != nil {
		// TODO: Look for module in search paths
		plg, err := plugin.Open(*pluginFile)
		if err != nil {
			fmt.Println("Plugin error: ", err)
			return
		}
		sym, err := plg.Lookup("Load")
		if err != nil {
			fmt.Println("Lookup error: ", err)
			return
		}
		scope = sym.(func() *object.Scope)()
	}

	if sourceFile := c.SearchSource(m); scope == nil && sourceFile != nil {
		text, _ := ioutil.ReadFile(*sourceFile)
		l := lexer.New(string(text))
		p := parser.New(l)
		program := p.ParseProgram()

		err := c.Compile(program)
		if err != nil {
			log.Fatalln("Compile error:", err)
		}
		return
	}

	if scope == nil {
		log.Fatalf("Failed to load scope: [%s], Exiting...\n", m)
	}

	si := c.addConstant(scope)

	c.emit(code.OpConstant, si)
	ssym := c.symbolTable.Define(m)
	c.saveSymbol(ssym)
	c.emit(code.OpPop)
}

func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(code.OpGetGlobal, s.Index)
	case LocalScope:
		c.emit(code.OpGetLocal, s.Index)
	case BuiltinScope:
		c.emit(code.OpGetBuiltin, s.Index)
	case FreeScope:
		c.emit(code.OpGetFree, s.Index)
	}
}

func (c *Compiler) saveSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(code.OpSetGlobal, s.Index)
	case LocalScope:
		c.emit(code.OpSetLocal, s.Index)
	case FreeScope:
		c.emit(code.OpSetFree, s.Index)
	}
}
