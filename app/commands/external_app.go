package commands

import (
	"fmt"
	"io"
	"os/exec"
)

func handleExternalApp(cmd string, args []string, in io.Reader, out io.Writer, errOut io.Writer) {
	_, err := exec.LookPath(cmd)
	if err != nil {
		fmt.Fprintf(errOut, "%s: command not found\n", cmd)
		return
	}
	executeExternalApp(cmd, args, in, out, errOut)
}

func executeExternalApp(app string, args []string, in io.Reader, out io.Writer, errOut io.Writer) {
	cmd := exec.Command(app, args...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = errOut
	cmd.Run()
}
