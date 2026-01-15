package commands

import (
	"slices"
	"strings"
)

func GetCompletions(prefix string) []string {
	var completions []string
	for cmd := range builtinCommands {
		if strings.HasPrefix(cmd, prefix) {
			completions = append(completions, cmd)
		}
	}

	executables := getPathFiles()
	for _, exe := range executables {
		parts := strings.Split(exe, "/")
		name := parts[len(parts)-1]
		if strings.HasPrefix(name, prefix) {
			completions = append(completions, name)
		}
	}

	slices.Sort(completions)
	return completions
}
