package repleval

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dreblang/core/evaluator"
	"github.com/dreblang/core/lexer"
	"github.com/dreblang/core/object"
	"github.com/dreblang/core/parser"
)

const Prompt = "evaluator /> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	globalEnv := object.NewEnvironment()

	for {
		fmt.Printf(Prompt)
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

		result := evaluator.Eval(program, globalEnv)
		if result != nil {
			io.WriteString(out, result.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
