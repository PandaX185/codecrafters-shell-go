package commands

import (
	"fmt"
	"os"
)

func HandlePwd() {
	if dir, err := os.Getwd(); err == nil {
		fmt.Println(dir)
	} else {
		fmt.Println("pwd: error retrieving current directory")
	}
}
