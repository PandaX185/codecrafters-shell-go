package commands

import (
	"fmt"
	"strings"
)

func HandleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}
