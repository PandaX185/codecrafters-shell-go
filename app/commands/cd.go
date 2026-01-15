package commands

import (
	"fmt"
	"os"
)

func HandleCd(dir string) (string, string) {
	if dir == "" || dir == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", "cd: Unable to determine home directory\n"
		}
		dir = homeDir
	}
	if err := os.Chdir(dir); err != nil {
		return "", fmt.Sprintf("cd: %v: No such file or directory\n", dir)
	}
	return "", ""
}
