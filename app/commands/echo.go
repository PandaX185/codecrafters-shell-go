package commands

import (
	"fmt"
	"strings"
)

func HandleEcho(args []string) {
	fmt.Println(strings.Trim(strings.Join(args, " "), "'"))
}
