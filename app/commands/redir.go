package commands

import (
	"io"
	"os"
)

type redir struct {
	OutFile   string
	OutAppend bool
	ErrFile   string
	ErrAppend bool
}

func parseRedirs(args []string) (redir, []string) {
	var r redir
	var clean []string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case ">", "1>":
			if i+1 < len(args) {
				r.OutFile = args[i+1]
				r.OutAppend = false
				i++
			}
		case ">>", "1>>":
			if i+1 < len(args) {
				r.OutFile = args[i+1]
				r.OutAppend = true
				i++
			}
		case "2>":
			if i+1 < len(args) {
				r.ErrFile = args[i+1]
				r.ErrAppend = false
				i++
			}
		case "2>>":
			if i+1 < len(args) {
				r.ErrFile = args[i+1]
				r.ErrAppend = true
				i++
			}
		default:
			clean = append(clean, args[i])
		}
	}
	return r, clean
}

func applyRedirs(out io.Writer, errOut io.Writer, r redir) (io.Writer, io.Writer, func()) {

	cleanup := func() {}
	if r.OutFile != "" {
		flags := os.O_CREATE | os.O_WRONLY
		if r.OutAppend {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}

		f, _ := os.OpenFile(r.OutFile, flags, 0644)
		out = f
		cleanup = func() {
			f.Close()
		}
	}

	if r.ErrFile != "" {
		flags := os.O_CREATE | os.O_WRONLY
		if r.ErrAppend {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}

		f, _ := os.OpenFile(r.ErrFile, flags, 0644)
		errOut = f
		prevCleanup := cleanup
		cleanup = func() {
			prevCleanup()
			f.Close()
		}
	}

	return out, errOut, cleanup
}
