package commands

import (
	"fmt"
)

func HandleType(cmd string) (string, string) {
	if _, ok := builtinCommands[cmd]; ok {
		return fmt.Sprintf("%s is a shell builtin\n", cmd), ""
	} else {
		path := pathSearch(cmd)
		if path != "" {
			return fmt.Sprintf("%s is %s\n", cmd, path), ""
		} else {
			return "", fmt.Sprintf("type: %s: not found\n", cmd)
		}
	}
}
