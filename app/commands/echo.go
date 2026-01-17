package commands

import (
	"fmt"
	"io"
	"strings"
)

func handleEcho(args []string, out io.Writer) {
	output := strings.Join(args, " ") + "\n"
	fmt.Fprint(out, output)
}
