package commands

import (
	"fmt"
	"strings"
)

func HandleEcho(args []string) {
	output := strings.Join(args, " ")
	output = unescapeString(output)
	fmt.Println(output)
}
