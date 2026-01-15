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
	out, err := exec.Command(app, args...).CombinedOutput()
	return string(out), err
}
