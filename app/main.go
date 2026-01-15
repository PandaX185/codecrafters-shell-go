package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/commands"
	"github.com/google/shlex"
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
		args, err := shlex.Split(cmd[len(cmdName):])
		if err != nil {
			fmt.Println("Error parsing arguments:", err)
			continue
		}
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
