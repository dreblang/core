package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/dreblang/core/compiler"
	"github.com/dreblang/core/lexer"
	"github.com/dreblang/core/object"
	"github.com/dreblang/core/parser"
	"github.com/dreblang/core/vm"
)

var input = `
let loopfunc = fn(n) {
	a = 0;
	sum = 0;
	loop (a < n) {
		sum = a + sum;
		a = a + 1;
	}
	return sum
}
loopfunc(50000000)
`

func main() {
	proffd, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(proffd)
	defer pprof.StopCPUProfile()

	flag.Parse()

	var duration time.Duration
	var result object.Object

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	comp := compiler.New()
	err := comp.Compile(program)
	if err != nil {
		fmt.Printf("compiler error: %s", err)
		return
	}

	machine := vm.New(comp.Bytecode())

	start := time.Now()

	err = machine.Run()
	if err != nil {
		fmt.Printf("vm error: %s", err)
		return
	}

	duration = time.Since(start)
	result = machine.LastPoppedStackElem()

	fmt.Printf(
		"engine=%s, result=%s, duration=%s\n",
		"vm",
		result.Inspect(),
		duration)
}
