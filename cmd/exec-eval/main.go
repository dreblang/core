package main

import (
	"io/ioutil"
	"os"

	"github.com/dreblang/core/evaluator"
	"github.com/dreblang/core/lexer"
	"github.com/dreblang/core/object"
	"github.com/dreblang/core/parser"
)

func main() {
	filename := os.Args[1]

	text, _ := ioutil.ReadFile(filename)
	l := lexer.New(string(text))
	p := parser.New(l)
	program := p.ParseProgram()
	globalEnv := object.NewEnvironment()
	evaluator.Eval(program, globalEnv)
}
