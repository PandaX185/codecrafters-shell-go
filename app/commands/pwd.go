package commands

import (
	"fmt"
	"io"
	"os"
)

func handlePwd(out io.Writer, errOut io.Writer) {
	if dir, err := os.Getwd(); err == nil {
		fmt.Fprintln(out, dir)
	} else {
		fmt.Fprintln(errOut, "pwd: error retrieving current directory")
	}
}
