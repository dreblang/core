package dreblc

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dreblang/core/compiler"
	"github.com/dreblang/core/lexer"
	"github.com/dreblang/core/object"
	"github.com/dreblang/core/parser"
	"github.com/dreblang/core/vm"
)

func Main() {
	filename := os.Args[1]

	text, _ := ioutil.ReadFile(filename)
	l := lexer.New(string(text))
	p := parser.New(l)
	program := p.ParseProgram()
	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalSize)
	symbolTable := compiler.NewSymbolTable()

	comp := compiler.NewWithState(symbolTable, constants)
	err := comp.Compile(program)
	if err != nil {
		fmt.Println("Compile error:", err)
		return
	}
	code := comp.Bytecode()
	constants = code.Constants

	machine := vm.NewWithGlobalsStore(code, globals)
	machine.Run()
}
