package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	for {
		fmt.Print("$ ")
		cmd, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading command:", err)
			continue
		}
		switch cmd[:len(cmd)-1] {
		case "exit":
			return
		default:
			fmt.Printf("%v: command not found\n", cmd[:len(cmd)-1])
		}
	}
}
