package replcompiler

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dreblang/core/compiler"
	"github.com/dreblang/core/lexer"
	"github.com/dreblang/core/object"
	"github.com/dreblang/core/parser"
	"github.com/dreblang/core/vm"
	"github.com/ttacon/chalk"
)

const Prompt = "jitc /> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalSize)
	symbolTable := compiler.NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	for {
		fmt.Printf("%s%s%s", chalk.Blue, Prompt, chalk.ResetColor)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			fmt.Printf("%s", chalk.Red)
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			fmt.Printf("%s", chalk.ResetColor)
			continue
		}

		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGlobalsStore(code, globals)
		err = machine.Run()
		if err != nil {
			fmt.Printf("%s", chalk.Red)
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			fmt.Printf("%s", chalk.ResetColor)
			continue
		}

		lastPopped := machine.LastPoppedStackElem()
		if lastPopped.Type() == object.ErrorObj {
			fmt.Printf("%s", chalk.Red)
		} else {
			fmt.Printf("%s", chalk.Green)
		}
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
		fmt.Printf("%s", chalk.ResetColor)
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
