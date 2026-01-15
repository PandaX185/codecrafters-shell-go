package commands

import (
	"fmt"
	"strings"

	"github.com/google/shlex"
)

func HandleEcho(args []string) {
	result, err := shlex.Split(strings.Join(args, " "))
	if err != nil {
		fmt.Println("echo: Error parsing arguments")
		return
	}
	fmt.Println(strings.Join(result, " "))
}
