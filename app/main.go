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
		if i, mode := commands.HasRedir(allArgs); i != -1 {
			args = allArgs[:i]
			fileName := allArgs[i+1]
			flag := os.O_CREATE | os.O_WRONLY
			if mode == 1 {
				flag |= os.O_APPEND
			} else {
				flag |= os.O_TRUNC
			}
			file, err := os.OpenFile(fileName, flag, 0644)
			if err != nil {
				fmt.Printf("Redirection error: %v\n", err)
				continue
			}
			defer file.Close()
			outFile = file
		}

		var (
			res    string
			errOut string
		)
		switch cmdName {
		case commands.Echo.String():
			res, errOut = commands.HandleEcho(args)
			break
		case commands.Type.String():
			cmd := strings.Join(args, " ")
			cmd = commands.UnescapeString(cmd)
			res, errOut = commands.HandleType(cmd)
			break
		case commands.Exit.String():
			return
		case commands.Pwd.String():
			res, errOut = commands.HandlePwd()
			break
		case commands.Cd.String():
			dir := strings.Join(args, " ")
			dir = commands.UnescapeString(dir)
			res, errOut = commands.HandleCd(dir)
			break
		default:
			res, errOut = commands.HandleExternalApp(cmdName, args)
		}
		outFile.WriteString(res)
		os.Stderr.WriteString(errOut)
	}
}
