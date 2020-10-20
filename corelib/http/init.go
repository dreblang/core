package http

import "github.com/dreblang/core/compiler"

func init() {
	compiler.RegisterLib("http", Load)
}
