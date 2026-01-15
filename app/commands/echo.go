package commands

import (
	"fmt"
	"strings"
)

func HandleEcho(args []string) (string, string) {
	return fmt.Sprint(strings.Join(args, " ")), ""
}
