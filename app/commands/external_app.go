package commands

import (
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
	out, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	errOut, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}
	outputBytes := make([]byte, 4096)
	n, _ := out.Read(outputBytes)
	if n > 0 {
		return string(outputBytes[:n]), nil
	}
	errBytes := make([]byte, 4096)
	n, _ = errOut.Read(errBytes)
	if n > 0 {
		return string(errBytes[:n]), nil
	}
	err = cmd.Wait()
	return string(outputBytes[:n]), err
}
