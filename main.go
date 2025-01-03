package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/MohamTahaB/interpreter-go/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! WELCOME TO THE MNKY CONSOLE !!!\n", user.Username)

	repl.Start(os.Stdin, os.Stdout)
}
