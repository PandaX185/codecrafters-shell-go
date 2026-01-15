package commands

import "strings"

func GetCompletions(prefix string) []string {
	var completions []string
	for cmd := range builtinCommands {
		if strings.HasPrefix(cmd, prefix) {
			completions = append(completions, cmd)
		}
	}
	return completions
}
