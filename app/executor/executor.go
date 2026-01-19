package executor

import (
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/commands"
)

func Execute(cmd string) {
	tokens := commands.Parse(cmd)
	if tokens == nil {
		return
	}

	cmdName := tokens[0]
	allArgs := tokens[1:]

	args, outFile, errFile, err := setupRedirections(allArgs)
	if err != nil {
		fmt.Printf("Redirection error: %v\n", err)
		return
	}
	defer func() {
		if outFile != os.Stdout {
			outFile.Close()
		}
		if errFile != os.Stderr {
			errFile.Close()
		}
	}()

	res, errOut := commands.ExecuteCommand(cmdName, args)

	if res != "" {
		res = strings.ReplaceAll(res, "\n", "\r\n")
		outFile.WriteString(res)
	}
	if errOut != "" {
		errOut = strings.ReplaceAll(errOut, "\n", "\r\n")
		errFile.WriteString(errOut)
	}
}

func setupRedirections(allArgs []string) (args []string, outFile, errFile *os.File, err error) {
	args = allArgs
	outFile = os.Stdout
	errFile = os.Stderr

	if i, outRedir := commands.HasOutRedir(allArgs); i != -1 {
		args = allArgs[:i]
		fileName := allArgs[i+1]
		flagOut := os.O_CREATE | os.O_WRONLY
		if outRedir == 1 {
			flagOut |= os.O_APPEND
		} else {
			flagOut |= os.O_TRUNC
		}
		outFile, err = os.OpenFile(fileName, flagOut, 0644)
		if err != nil {
			return
		}
	}

	if i, errRedir := commands.HasErrRedir(allArgs); i != -1 {
		args = allArgs[:min(i, len(args))]
		fileName := allArgs[i+1]
		flagErr := os.O_CREATE | os.O_WRONLY
		if errRedir == 1 {
			flagErr |= os.O_APPEND
		} else {
			flagErr |= os.O_TRUNC
		}
		errFile, err = os.OpenFile(fileName, flagErr, 0644)
		if err != nil {
			return
		}
	}

	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
