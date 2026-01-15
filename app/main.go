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
			cmd := strings.Join(args, " ")
			cmd = commands.UnescapeString(cmd)
			commands.HandleType(cmd)
			break
		case commands.Exit.String():
			return
		case commands.Pwd.String():
			commands.HandlePwd()
			break
		case commands.Cd.String():
			dir := strings.Join(args, " ")
			dir = commands.UnescapeString(dir)
			commands.HandleCd(dir)
			break
		default:
			commands.HandleExternalApp(cmdName, args)
		}
	}
}
