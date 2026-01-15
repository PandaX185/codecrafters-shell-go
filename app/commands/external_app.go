package commands

import (
	"fmt"
	"os/exec"
)

func HandleExternalApp(cmd string, args []string) {
	apps := getPathFiles()
	for _, app := range apps {
		if getFileName(app) == cmd {
			if err := executeExternalApp(app, args); err != nil {
				fmt.Printf("Error executing %s: %v\n", cmd, err)
			}
			fmt.Printf("Program was passed %v args (including program name)\n", len(args)+1)
			return
		}
	}
	fmt.Printf("%s: command not found\n", cmd)
}

func executeExternalApp(app string, args []string) error {
	if err := exec.Command(app, args...).Run(); err != nil {
		return err
	}
	return nil
}
