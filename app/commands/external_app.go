package commands

import (
	"bytes"
	"fmt"
	"os/exec"
)

func HandleExternalApp(cmd string, args []string) {
	if output, err := executeExternalApp(cmd, args); err == nil {
		fmt.Print(output)
		return
	}
	fmt.Printf("%s: command not found\n", cmd)
}

func executeExternalApp(app string, args []string) (string, error) {
	cmd := exec.Command(app, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	output := stdout.String()
	errOutput := stderr.String()
	if err != nil {
		return output + errOutput, err
	}
	return output, nil
}
