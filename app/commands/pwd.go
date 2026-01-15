package commands

import (
	"os"
)

func HandlePwd() string {
	if dir, err := os.Getwd(); err == nil {
		return dir + "\n"
	} else {
		return "pwd: error retrieving current directory\n"
	}
}
