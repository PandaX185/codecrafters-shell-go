package commands

import "os"

func HandleRedir(result string, outFile *os.File) {
	outFile.WriteString(result)
}

func HasRedir(args []string) int {
	for i, arg := range args {
		if arg == ">" {
			return i
		}
	}
	return -1
}
