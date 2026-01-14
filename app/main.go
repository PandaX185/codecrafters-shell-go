package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
		case "echo":
			fmt.Println(strings.Join(args, " "))
			break
		case "exit":
			return
		default:
			fmt.Printf("%v: command not found\n", cmdName)
		}
	}
}
