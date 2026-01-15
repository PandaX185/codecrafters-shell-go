package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/commands"
	"golang.org/x/term"
)

func main() {
	fd := int(os.Stdin.Fd())
	old, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("Error setting terminal to raw mode:", err)
		return
	}
	defer term.Restore(fd, old)

	for {
		fmt.Print("$ ")

		var cmd string
		consoleReader := bufio.NewReader(os.Stdin)
		for {
			b, err := consoleReader.ReadByte()
			if err != nil {
				fmt.Println("Error reading from console:", err)
				return
			}
			if b == '\n' || b == '\r' {
				fmt.Print("\r\n")
				break
			}
			if b == '\t' {
				completions := commands.GetCompletions(cmd)
				if len(completions) == 1 {
					toAdd := completions[0][len(cmd):] + " "
					cmd += toAdd
					fmt.Print(toAdd)
				} else if len(completions) > 1 {
					fmt.Print("\r\n")
					for _, c := range completions {
						fmt.Printf("%s\t", c)
					}
					fmt.Print("\r\n$ " + cmd)
				}
				continue
			}
			if b == 3 {
				fmt.Print("^C\r\n")
				cmd = ""
				return
			}
			if b == 127 || b == 8 {
				if len(cmd) > 0 {
					cmd = cmd[:len(cmd)-1]
					fmt.Print("\b \b")
				}
				continue
			}

			cmd += string(b)
			fmt.Printf("%c", b)
		}

		tokens := commands.Parse(cmd)
		if tokens == nil {
			continue
		}

		cmdName := tokens[0]
		allArgs := tokens[1:]
		args := allArgs
		outFile := os.Stdout
		errFile := os.Stderr
		if i, outRedir := commands.HasOutRedir(allArgs); i != -1 {
			args = allArgs[:i]
			fileName := allArgs[i+1]
			flagOut := os.O_CREATE | os.O_WRONLY
			if outRedir == 1 {
				flagOut |= os.O_APPEND
			} else {
				flagOut |= os.O_TRUNC
			}
			file, err := os.OpenFile(fileName, flagOut, 0644)
			if err != nil {
				fmt.Printf("Output redirection error: %v\n", err)
				continue
			}
			defer file.Close()
			outFile = file
		}
		if i, errRedir := commands.HasErrRedir(allArgs); i != -1 {
			args = allArgs[:min(i, len(args))]
			fileName := allArgs[i+1]
			flagErr := os.O_CREATE | os.O_WRONLY
			if errRedir == 1 {
				flagErr |= os.O_APPEND
			} else {
				flagErr |= os.O_TRUNC
			}
			file, err := os.OpenFile(fileName, flagErr, 0644)
			if err != nil {
				fmt.Printf("Error redirection error: %v\n", err)
				continue
			}
			defer file.Close()
			errFile = file
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
		if res != "" {
			outFile.WriteString(res + "\r\n")
		}
		if errOut != "" {
			errFile.WriteString(errOut + "\r\n")
		}
	}
}
