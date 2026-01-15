package commands

import (
	"fmt"
	"os"
)

func HandleCd(dir string) {
	if dir == "" || dir == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("cd: Unable to determine home directory")
			return
		}
		dir = homeDir
	}
	if err := os.Chdir(dir); err != nil {
		fmt.Printf("cd: %v: No such file or directory\n", dir)
	}
}
