package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/dreblang/core/repleval"
)

func main() {
	user, err := user.Current()
	if err != nil {
		print(err)
	}

	fmt.Printf("Hello %s! This is the Dreblang!\n", user.Username)
	fmt.Printf("Waiting for your commands...\n")
	repleval.Start(os.Stdin, os.Stdout)
}
