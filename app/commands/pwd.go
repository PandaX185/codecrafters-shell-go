package commands

import (
	"os"
)

func handlePwd() (string, string) {
	if dir, err := os.Getwd(); err == nil {
		return dir + "\n", ""
	} else {
		return "", "pwd: error retrieving current directory"
	}
}
