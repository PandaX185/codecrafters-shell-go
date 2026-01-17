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
		tabCount := 0
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
				tabCount = tabCount%2 + 1
				completions := commands.GetCompletions(cmd)
				if len(completions) == 1 {
					toAdd := completions[0] + " "
					fmt.Print(toAdd[len(cmd):])
					cmd = toAdd
				} else {
					if tabCount == 1 {
						toAdd := commands.GetLcsPrefix(completions)
						if len(completions) > 0 && toAdd != cmd {
							toAdd = toAdd[len(cmd):]
							fmt.Print(toAdd)
							cmd += toAdd
						} else {
							fmt.Printf("%c", 0x07)
						}
					} else {
						fmt.Print("\r\n")
						fmt.Print(strings.Join(completions, "  ") + "\r\n")
						fmt.Print("$ " + cmd)
					}
				}
				continue
			}
			if b == 3 {
				fmt.Print("^C\r\n")
				cmd = ""
				return
			}
			if b == 127 || b == 8 {
				tabCount = 0
				if len(cmd) > 0 {
					cmd = cmd[:len(cmd)-1]
					fmt.Print("\b \b")
				}
				continue
			}

			tabCount = 0
			cmd += string(b)
			fmt.Printf("%c", b)
		}

		tokens := commands.Parse(cmd)
		if len(tokens) == 0 {
			continue
		}

		commands.HandlePipeline(commands.ParsePipeline(cmd))
	}
}
