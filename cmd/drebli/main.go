package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/dreblang/core/replcompiler"
)

func main() {
	user, err := user.Current()
	if err != nil {
		print(err)
	}

	fmt.Printf("Hello %s! This is the Dreblang!\n", user.Username)
	fmt.Printf("Waiting for your commands...\n")
	replcompiler.Start(os.Stdin, os.Stdout)
}
