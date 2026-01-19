package commands

import (
	"strings"
)

func handleEcho(args []string) (string, string) {
	return strings.Join(args, " ") + "\n", ""
}
