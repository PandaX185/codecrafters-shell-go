package commands

import (
	"fmt"
	"strings"
)

func HandleEcho(args []string) string {
	return fmt.Sprintln(strings.Join(args, " "))
}
