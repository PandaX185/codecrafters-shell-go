package commands

import (
	"fmt"
	"os"
	"strings"
)

func HandleType(cmd string) {
	if _, ok := builtinCommands[cmd]; ok {
		fmt.Printf("%s is a shell builtin\n", cmd)
	} else if !pathSearch(cmd) {
		fmt.Printf("%s: not found\n", cmd)
	}
}

func pathSearch(cmd string) bool {
	path := os.Getenv("PATH")
	dirs := strings.Split(path, ":")
	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			if file.Name() == cmd {
				info, err := file.Info()
				if err != nil {
					continue
				}
				if checkFileExecutable(info) {
					fmt.Printf("%s is %s/%s\n", cmd, dir, cmd)
					return true
				}
			}
		}

	}
	return false
}

func checkFileExecutable(file os.FileInfo) bool {
	return file.Mode().IsRegular() && (file.Mode().Perm()&0111 != 0)
}
