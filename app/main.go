package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/commands"
)

func main() {
	for {
		fmt.Print("$ ")
		cmd, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command:", err)
			continue
		}
		cmd = cmd[:len(cmd)-1]
		cmdName := strings.Split(cmd, " ")[0]
		args := strings.Split(cmd, " ")[1:]
		switch cmdName {
		case commands.Echo.String():
			commands.HandleEcho(args)
			break
		case commands.Type.String():
			commands.HandleType(strings.Join(args, " "))
			break
		case commands.Exit.String():
			return
		default:
			fmt.Printf("%v: command not found\n", cmdName)
		}
	}
}
