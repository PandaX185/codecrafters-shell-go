package commands

import (
	"fmt"
	"io"
)

func handleType(cmd string, out io.Writer, errOut io.Writer) {
	if _, ok := builtinCommands[cmd]; ok {
		fmt.Fprintf(out, "%s is a shell builtin\n", cmd)
	} else {
		path := pathSearch(cmd)
		if path != "" {
			fmt.Fprintf(out, "%s is %s\n", cmd, path)
		} else {
			fmt.Fprintf(errOut, "%s: not found\n", cmd)
		}
	}
}
