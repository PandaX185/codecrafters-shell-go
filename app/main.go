package main

import (
	"github.com/codecrafters-io/shell-starter-go/app/shell"
)

func main() {
	shell.SetupTerminal()
	defer shell.RestoreTerminal()

	shell.Repl()
}
