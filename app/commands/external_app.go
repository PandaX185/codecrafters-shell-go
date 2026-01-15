package commands

import (
	"bytes"
	"fmt"
	"os/exec"
)

func HandleExternalApp(cmd string, args []string) (string, string) {
	_, err := exec.LookPath(cmd)
	if err != nil {
		return "", fmt.Sprintf("%s: command not found\n", cmd)
	}
	return executeExternalApp(cmd, args)
}

func executeExternalApp(app string, args []string) (string, string) {
	cmd := exec.Command(app, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()

	return stdout.String(), stderr.String()
}
