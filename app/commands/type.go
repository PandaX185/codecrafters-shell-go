package commands

import "fmt"

func HandleType(cmd string) {
	if _, ok := builtinCommands[cmd]; ok {
		fmt.Printf("%s is a shell builtin\n", cmd)
	} else {
		fmt.Printf("%s: not found\n", cmd)
	}
}
