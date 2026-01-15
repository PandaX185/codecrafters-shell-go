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
		allArgs := tokens[1:]
		args := allArgs
		outFile := os.Stdout
		if i := commands.HasRedir(allArgs); i != -1 {
			args = allArgs[:i]
			fileName := allArgs[i+1]
			file, err := os.Create(fileName)
			if err != nil {
				fmt.Printf("Redirection error: %v\n", err)
				continue
			}
			defer file.Close()
			outFile = file
		}

		var res string
		switch cmdName {
		case commands.Echo.String():
			res = commands.HandleEcho(args)
			break
		case commands.Type.String():
			cmd := strings.Join(args, " ")
			cmd = commands.UnescapeString(cmd)
			res = commands.HandleType(cmd)
			break
		case commands.Exit.String():
			return
		case commands.Pwd.String():
			res = commands.HandlePwd()
			break
		case commands.Cd.String():
			dir := strings.Join(args, " ")
			dir = commands.UnescapeString(dir)
			res = commands.HandleCd(dir)
			break
		default:
			res = commands.HandleExternalApp(cmdName, args)
		}
		commands.HandleRedir(res, outFile)
	}
}
