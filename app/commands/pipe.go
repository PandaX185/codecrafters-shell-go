package commands

import (
	"io"
	"os"
	"strings"
	"sync"

	"golang.org/x/term"
)

type PipeCmd struct {
	Cmd  string
	Args []string
}

func ParsePipeline(line string) []PipeCmd {
	pipes := strings.Split(line, "|")
	pipeline := []PipeCmd{}
	for _, pipeStr := range pipes {
		pipeTokens := Parse(strings.TrimSpace(pipeStr))
		pipeCmd := ParsePipe(pipeTokens)
		pipeline = append(pipeline, pipeCmd)
	}
	return pipeline
}

func ParsePipe(args []string) PipeCmd {
	pipe := PipeCmd{}
	if len(args) == 0 {
		return pipe
	}

	pipe.Cmd = args[0]
	if len(args) > 1 {
		pipe.Args = args[1:]
	}

	return pipe
}

func HandlePipeline(pipes []PipeCmd) {
	var prevReader io.Reader = os.Stdin
	var wg sync.WaitGroup
	var terminalOut io.Writer = os.Stdout
	var terminalErr io.Writer = os.Stderr

	if term.IsTerminal(int(os.Stdout.Fd())) {
		terminalOut = &crlfWriter{w: os.Stdout}
	}
	if term.IsTerminal(int(os.Stderr.Fd())) {
		terminalErr = &crlfWriter{w: os.Stderr}
	}

	for i, pipe := range pipes {
		isLast := i == len(pipes)-1

		var out io.Writer
		var nextReader io.Reader

		if isLast {
			out = terminalOut
		} else {
			pr, pw := io.Pipe()
			out = pw
			nextReader = pr
		}

		wg.Add(1)
		go func(cmd string, args []string, in io.Reader, out io.Writer, isLast bool) {
			defer wg.Done()
			ExecuteCommand(cmd, args, in, out, terminalErr)
			if pw, ok := out.(*io.PipeWriter); ok {
				pw.Close()
			}
		}(pipe.Cmd, pipe.Args, prevReader, out, isLast)
		if !isLast {
			prevReader = nextReader
		}
	}
	wg.Wait()
}
