package commands

import (
	"bufio"
	"context"
	"io"
	"os"
	"strings"
	"sync"
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

func HandlePipeline(ctx context.Context, cancel context.CancelFunc, pipes []PipeCmd, stdout *bufio.Writer, stderr io.Writer) {
	var prevReader io.Reader = os.Stdin
	var wg sync.WaitGroup

	go func() {
		buf := make([]byte, 1)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := os.Stdin.Read(buf)
				if err != nil {
					return
				}
				if n == 1 && buf[0] == 3 {
					cancel()
					return
				}
			}
		}
	}()

	for i, pipe := range pipes {
		isLast := i == len(pipes)-1

		var out io.Writer
		var nextReader io.Reader

		if isLast {
			go func(r io.Reader, w *bufio.Writer) {
				buf := make([]byte, 1024)
				for {
					n, err := r.Read(buf)
					if err != nil || n == 0 {
						break
					}
					for i := range n {
						if buf[i] == '\n' {
							w.Write([]byte{'\r', '\n'})
						} else {
							w.Write([]byte{buf[i]})
						}
						w.Flush()
					}
				}
			}(prevReader, stdout)
		} else {
			pr, pw := io.Pipe()
			out = pw
			nextReader = pr
		}

		wg.Add(1)
		go func(ctx context.Context, cmd string, args []string, in io.Reader, out io.Writer, isLast bool) {
			defer wg.Done()
			ExecuteCommand(ctx, cmd, args, in, out, stderr)
			if pw, ok := out.(*io.PipeWriter); ok {
				pw.Close()
			}
		}(ctx, pipe.Cmd, pipe.Args, prevReader, out, isLast)
		if !isLast {
			prevReader = nextReader
		}
	}
	wg.Wait()
}
