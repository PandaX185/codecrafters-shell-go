package commands

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

func handleExternalApp(ctx context.Context, cmd string, args []string, in io.Reader, out io.Writer, errOut io.Writer) {
	_, err := exec.LookPath(cmd)
	if err != nil {
		fmt.Fprintf(errOut, "%s: command not found\n", cmd)
		return
	}
	executeExternalApp(ctx, cmd, args, in, out, errOut)
}

func executeExternalApp(ctx context.Context, app string, args []string, in io.Reader, out io.Writer, errOut io.Writer) {
	cmd := exec.CommandContext(ctx, app, args...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = errOut
	cmd.Run()
}
