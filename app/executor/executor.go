package executor

import (
	"fmt"
	"io"
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

		if res != "\n" {
			res = strings.ReplaceAll(res, "\n", "\r\n")
			outFile.WriteString(res)
		}
		if errOut != "\n" {
			errOut = strings.ReplaceAll(errOut, "\n", "\r\n")
			errFile.WriteString(errOut)
		}
	} else {
		executePipeline(pipeline)
	}
}

type builtinCmd struct {
	name   string
	args   []string
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (b *builtinCmd) Run() error {
	var input string
	if b.stdin != nil && b.stdin != os.Stdin {
		data, _ := io.ReadAll(b.stdin)
		input = string(data)
	}

	args := b.args
	if input != "" && b.name == "cat" {
		b.stdout.Write([]byte(input))
		return nil
	}

	res, errOut := commands.ExecuteCommand(b.name, args)
	res += "\n"
	errOut += "\n"

	if res != "\n" {
		b.stdout.Write([]byte(res))
	}
	if errOut != "\n" {
		b.stderr.Write([]byte(errOut))
	}
	return nil
}

func executePipeline(pipeline [][]string) {
	type pipelineCmd struct {
		isBuiltin bool
		builtin   *builtinCmd
		external  *exec.Cmd
	}

	cmds := make([]pipelineCmd, len(pipeline))
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

		var stdin io.Reader = inFile
		if inFile == os.Stdin && i > 0 {
			stdin = pipeReaders[i-1]
		} else if inFile != os.Stdin {
			filesToClose = append(filesToClose, inFile)
		}
		var stdout io.Writer = outFile
		if i < len(pipeline)-1 {
			stdout = pipeWriters[i]
		} else if outFile != os.Stdout {
			filesToClose = append(filesToClose, outFile)
		}
		var stderr io.Writer = errFile
		if errFile != os.Stderr {
			filesToClose = append(filesToClose, errFile)
		}

		if commands.IsBuiltin(cmdName) {
			cmds[i] = pipelineCmd{
				isBuiltin: true,
				builtin: &builtinCmd{
					name:   cmdName,
					args:   args,
					stdin:  stdin,
					stdout: stdout,
					stderr: stderr,
				},
			}
		} else {
			cmd := exec.Command(cmdName, args...)
			if f, ok := stdin.(*os.File); ok {
				cmd.Stdin = f
			}
			if f, ok := stdout.(*os.File); ok {
				cmd.Stdout = f
			} else {
				cmd.Stdout = stdout
			}
			if f, ok := stderr.(*os.File); ok {
				cmd.Stderr = f
			} else {
				cmd.Stderr = stderr
			}
			cmds[i] = pipelineCmd{
				isBuiltin: false,
				external:  cmd,
			}
		}
	}
	for i := 0; i < len(cmds); i++ {
		if !cmds[i].isBuiltin && cmds[i].external != nil {
			err := cmds[i].external.Start()
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
		if cmds[i].isBuiltin && cmds[i].builtin != nil {
			cmds[i].builtin.Run()
		} else if cmds[i].external != nil {
			cmds[i].external.Wait()
		}
	}
	for _, r := range pipeReaders {
		r.Close()
	}
	for _, f := range filesToClose {
		f.Close()
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
