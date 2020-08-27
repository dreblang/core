package main

import (
	"io/ioutil"
	"os"

	"github.com/dreblang/core/compiler"
	"github.com/dreblang/core/lexer"
	"github.com/dreblang/core/object"
	"github.com/dreblang/core/parser"
	"github.com/dreblang/core/vm"
)

func main() {
	filename := os.Args[1]

	text, _ := ioutil.ReadFile(filename)
	l := lexer.New(string(text))
	p := parser.New(l)
	program := p.ParseProgram()
	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalSize)
	symbolTable := compiler.NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	comp := compiler.NewWithState(symbolTable, constants)
	comp.Compile(program)
	code := comp.Bytecode()
	constants = code.Constants

	machine := vm.NewWithGlobalsStore(code, globals)
	machine.Run()
}
