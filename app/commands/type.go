package commands

import (
	"fmt"
)

func handleType(cmd string) (string, string) {
	if _, ok := builtinCommands[cmd]; ok {
		return fmt.Sprintf("%s is a shell builtin", cmd), ""
	} else {
		path := pathSearch(cmd)
		if path != "" {
			return fmt.Sprintf("%s is %s", cmd, path), ""
		} else {
			return "", fmt.Sprintf("%s: not found", cmd)
		}
	}
}
