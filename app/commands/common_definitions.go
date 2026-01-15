package commands

import (
	"fmt"
	"os"
	"strings"
)

type command int

const (
	Echo command = iota
	Exit
	Type
	Pwd
	Cd
)

var builtinCommands = map[string]bool{
	Echo.String(): true,
	Exit.String(): true,
	Type.String(): true,
	Pwd.String():  true,
	Cd.String():   true,
}

func (c command) String() string {
	switch c {
	case Echo:
		return "echo"
	case Exit:
		return "exit"
	case Type:
		return "type"
	case Pwd:
		return "pwd"
	case Cd:
		return "cd"
	default:
		return "unknown"
	}
}

func getPathFiles() (result []string) {
	path := os.Getenv("PATH")
	dirs := strings.Split(path, ":")

	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			info, err := file.Info()
			if err != nil {
				continue
			}
			if checkFileExecutable(info) {
				result = append(result, dir+"/"+file.Name())
			}
		}
	}

	return result
}

func pathSearch(cmd string) bool {
	files := getPathFiles()
	for _, file := range files {
		if strings.HasSuffix(file, "/"+cmd) {
			fmt.Println(file)
			return true
		}
	}
	return false
}

func checkFileExecutable(file os.FileInfo) bool {
	return file.Mode().Perm()&0111 != 0
}

func getFileName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
