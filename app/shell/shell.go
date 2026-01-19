package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/commands"
	"github.com/codecrafters-io/shell-starter-go/app/executor"
	"golang.org/x/term"
)

var oldState *term.State

func SetupTerminal() {
	fd := int(os.Stdin.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("Error setting terminal to raw mode:", err)
		os.Exit(1)
	}
	oldState = state
}

func RestoreTerminal() {
	if oldState != nil {
		fd := int(os.Stdin.Fd())
		term.Restore(fd, oldState)
	}
}

func Repl() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		cmd, err := ReadCommand(reader)
		if err != nil {
			fmt.Println("Error reading command:", err)
			continue
		}

		if cmd == "" {
			continue
		}

		RestoreTerminal()
		executor.Execute(cmd)
		SetupTerminal()
	}
}

func ReadCommand(reader *bufio.Reader) (string, error) {
	var cmd string
	tabCount := 0
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return "", err
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
						fmt.Printf("%c", 0x07) // Bell
					}
				} else {
					fmt.Print("\r\n")
					fmt.Print(strings.Join(completions, "  ") + "\r\n")
					fmt.Print("$ " + cmd)
				}
			}
			continue
		}
		if b == 3 { // Ctrl+C
			fmt.Print("^C\r\n")
			return "", nil
		}
		if b == 127 || b == 8 { // Backspace
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
	return cmd, nil
}
