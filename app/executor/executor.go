package executor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/commands"
)

func Execute(cmd string) {
	pipeline := commands.Parse(cmd)
	if pipeline == nil || len(pipeline) == 0 {
		return
	}

	if len(pipeline) == 1 {
		tokens := pipeline[0]
		if len(tokens) == 0 {
			return
		}

		cmdName := tokens[0]
		allArgs := tokens[1:]

		args, inFile, outFile, errFile, err := setupRedirections(allArgs)
		if err != nil {
			fmt.Printf("Redirection error: %v\n", err)
			return
		}
		defer func() {
			if inFile != os.Stdin {
				inFile.Close()
			}
			if outFile != os.Stdout {
				outFile.Close()
			}
			if errFile != os.Stderr {
				errFile.Close()
			}
		}()

		res, errOut := commands.ExecuteCommand(cmdName, args)
		res += "\n"
		errOut += "\n"

		if res != "" {
			res = strings.ReplaceAll(res, "\n", "\r\n")
			outFile.WriteString(res)
		}
		if errOut != "" {
			errOut = strings.ReplaceAll(errOut, "\n", "\r\n")
			errFile.WriteString(errOut)
		}
	} else {

		cmds := make([]*exec.Cmd, len(pipeline))
		var pipeReaders, pipeWriters []*os.File
		var filesToClose []*os.File

		for i := 0; i < len(pipeline)-1; i++ {
			r, w, err := os.Pipe()
			if err != nil {
				fmt.Printf("Pipe error: %v\n", err)
				return
			}
			pipeReaders = append(pipeReaders, r)
			pipeWriters = append(pipeWriters, w)
		}

		for i, tokens := range pipeline {
			if len(tokens) == 0 {
				continue
			}

			cmdName := tokens[0]
			allArgs := tokens[1:]

			args, inFile, outFile, errFile, err := setupRedirections(allArgs)
			if err != nil {
				fmt.Printf("Redirection error: %v\n", err)
				return
			}

			cmds[i] = exec.Command(cmdName, args...)

			if inFile != os.Stdin {
				cmds[i].Stdin = inFile
				filesToClose = append(filesToClose, inFile)
			} else if i > 0 {
				cmds[i].Stdin = pipeReaders[i-1]
			}

			if i < len(pipeline)-1 {
				cmds[i].Stdout = pipeWriters[i]
			} else {
				if outFile != os.Stdout {
					cmds[i].Stdout = outFile
					filesToClose = append(filesToClose, outFile)
				} else {
					cmds[i].Stdout = os.Stdout
				}
			}

			if errFile != os.Stderr {
				cmds[i].Stderr = errFile
				filesToClose = append(filesToClose, errFile)
			} else {
				cmds[i].Stderr = os.Stderr
			}
		}

		for i := 0; i < len(cmds); i++ {
			if cmds[i] != nil {
				err := cmds[i].Start()
				if err != nil {
					fmt.Printf("Start error: %v\n", err)
					return
				}
			}
		}

		for _, w := range pipeWriters {
			w.Close()
		}

		for i := 0; i < len(cmds); i++ {
			if cmds[i] != nil {
				cmds[i].Wait()
			}
		}

		for _, r := range pipeReaders {
			r.Close()
		}

		for _, f := range filesToClose {
			f.Close()
		}
	}
}

func setupRedirections(allArgs []string) (args []string, inFile, outFile, errFile *os.File, err error) {
	args = make([]string, 0, len(allArgs))
	inFile = os.Stdin
	outFile = os.Stdout
	errFile = os.Stderr

	i := 0
	for i < len(allArgs) {
		arg := allArgs[i]
		if arg == "<" {
			if i+1 < len(allArgs) {
				fileName := allArgs[i+1]
				inFile, err = os.Open(fileName)
				if err != nil {
					return
				}
				i += 2
				continue
			}
		} else if arg == "<<" {

			i += 2
			continue
		} else if arg == ">" || arg == "1>" {
			if i+1 < len(allArgs) {
				fileName := allArgs[i+1]
				flag := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
				outFile, err = os.OpenFile(fileName, flag, 0644)
				if err != nil {
					return
				}
				i += 2
				continue
			}
		} else if arg == ">>" || arg == "1>>" {
			if i+1 < len(allArgs) {
				fileName := allArgs[i+1]
				flag := os.O_CREATE | os.O_WRONLY | os.O_APPEND
				outFile, err = os.OpenFile(fileName, flag, 0644)
				if err != nil {
					return
				}
				i += 2
				continue
			}
		} else if arg == "2>" {
			if i+1 < len(allArgs) {
				fileName := allArgs[i+1]
				flag := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
				errFile, err = os.OpenFile(fileName, flag, 0644)
				if err != nil {
					return
				}
				i += 2
				continue
			}
		} else if arg == "2>>" {
			if i+1 < len(allArgs) {
				fileName := allArgs[i+1]
				flag := os.O_CREATE | os.O_WRONLY | os.O_APPEND
				errFile, err = os.OpenFile(fileName, flag, 0644)
				if err != nil {
					return
				}
				i += 2
				continue
			}
		}
		args = append(args, arg)
		i++
	}
	return
}
