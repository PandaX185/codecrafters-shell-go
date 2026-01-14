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
		fmt.Printf("%v: command not found\n", cmd[:len(cmd)-1])
	}
}
