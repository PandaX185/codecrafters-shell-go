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

		tokens := commands.Parse(cmd)
		if tokens == nil {
			continue
		}

		cmdName := tokens[0]
		args := tokens[1:]
		switch cmdName {
		case commands.Echo.String():
			commands.HandleEcho(args)
			break
		case commands.Type.String():
			commands.HandleType(strings.Join(args, " "))
			break
		case commands.Exit.String():
			return
		case commands.Pwd.String():
			commands.HandlePwd()
			break
		case commands.Cd.String():
			commands.HandleCd(strings.Join(args, " "))
			break
		default:
			commands.HandleExternalApp(cmdName, args)
		}
	}
}
