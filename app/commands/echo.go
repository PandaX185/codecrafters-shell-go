package commands

import (
	"fmt"
	"strings"
)

func HandleEcho(args []string) {
	output := strings.Join(args, " ")
	output = UnescapeString(output)
	fmt.Println(output)
}
