package commands

import (
	"bytes"
	"fmt"
	"os/exec"
)

func HandleExternalApp(cmd string, args []string) {
	_, err := exec.LookPath(cmd)
	if err != nil {
		fmt.Printf("%s: command not found\n", cmd)
		return
	}
	output, errOutput := executeExternalApp(cmd, args)
	fmt.Print(output)
	fmt.Print(errOutput)
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
