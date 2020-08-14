package repleval

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ttacon/chalk"

	"github.com/dreblang/core/evaluator"
	"github.com/dreblang/core/lexer"
	"github.com/dreblang/core/object"
	"github.com/dreblang/core/parser"
)

const Prompt = "eval > "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	globalEnv := object.NewEnvironment()

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

		result := evaluator.Eval(program, globalEnv)
		if result != nil {
			if result.Type() == object.ErrorObj {
				fmt.Printf("%s", chalk.Red)
			} else {
				fmt.Printf("%s", chalk.Green)
			}
			io.WriteString(out, result.Inspect())
			io.WriteString(out, "\n")
			fmt.Printf("%s", chalk.ResetColor)
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
