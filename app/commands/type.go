package commands

import (
	"fmt"
)

func HandleType(cmd string) string {
	if _, ok := builtinCommands[cmd]; ok {
		return fmt.Sprintf("%s is a shell builtin\n", cmd)
	} else if !pathSearch(cmd) {
		return fmt.Sprintf("%s: not found\n", cmd)
	}
	return ""
}
