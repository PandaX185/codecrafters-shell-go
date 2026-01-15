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

	slices.SortFunc(completions, func(i, j string) int {
		if len(i) != len(j) {
			return len(i) - len(j)
		}
		return strings.Compare(i, j)
	})
	completions = slices.Compact(completions)
	return completions
}

func GetLcsPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	result := ""
	minLen := len(strs[0])

	for i := 0; i < minLen; i++ {
		c := strs[0][i]
		for j := 1; j < len(strs); j++ {
			if strs[j][i] != c {
				return result
			}
		}
		result += string(c)
	}
	return result
}
