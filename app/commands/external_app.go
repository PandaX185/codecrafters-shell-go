package commands

import (
	"fmt"
	"os/exec"
)

func HandleExternalApp(cmd string, args []string) {
	apps := getPathFiles()
	for _, app := range apps {
		if getFileName(app) == cmd {
			if out, err := executeExternalApp(app, args); err != nil {
				fmt.Printf("Error executing %s: %v\n", cmd, err)
			} else {
				fmt.Print(out)
			}
			return
		}
	}
	fmt.Printf("%s: command not found\n", cmd)
}

func executeExternalApp(app string, args []string) (string, error) {
	out, err := exec.Command(app, args...).CombinedOutput()
	return string(out), err
}
