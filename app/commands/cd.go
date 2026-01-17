package commands

import (
	"fmt"
	"io"
	"os"
)

func handleCd(dir string, out io.Writer, errOut io.Writer) {
	if dir == "" || dir == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(errOut, "cd: Unable to determine home directory")
			return
		}
		dir = homeDir
	}
	if err := os.Chdir(dir); err != nil {
		fmt.Fprintf(errOut, "cd: %v: No such file or directory\n", dir)
	}
}
