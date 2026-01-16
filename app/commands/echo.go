package commands

import (
	"fmt"
	"strings"
)

func handleEcho(args []string) (string, string) {
	return fmt.Sprintln(strings.Join(args, " ")), ""
}
