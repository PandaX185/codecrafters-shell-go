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

func unescapeString(s string) string {
	runes := []rune(s)
	result := make([]rune, 0, len(runes))
	i := 0
	for i < len(runes) {
		if runes[i] == '\\' && i+1 < len(runes) && isOctalDigit(runes[i+1]) {
			j := i + 1
			for j < len(runes) && j < i+4 && isOctalDigit(runes[j]) {
				j++
			}
			octalStr := string(runes[i+1 : j])
			if val, err := parseOctal(octalStr); err == nil {
				result = append(result, rune(val))
				i = j
				continue
			}
		}
		result = append(result, runes[i])
		i++
	}
	s = string(result)

	s = strings.ReplaceAll(s, `\\`, `\`)
	s = strings.ReplaceAll(s, `\n`, "\n")
	s = strings.ReplaceAll(s, `\t`, "\t")
	s = strings.ReplaceAll(s, `\"`, `"`)
	s = strings.ReplaceAll(s, `\'`, "'")
	return s
}

func isOctalDigit(r rune) bool {
	return r >= '0' && r <= '7'
}

func parseOctal(s string) (int, error) {
	val := 0
	for _, r := range s {
		val = val*8 + int(r-'0')
	}
	return val, nil
}
