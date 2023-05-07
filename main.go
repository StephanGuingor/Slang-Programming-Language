package main

import (
	"compiler-book/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	// args := os.Args
	
	if len(os.Args) == 2 {
		repl.StartFile(os.Args[1])
		return
	}

	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Slang programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
