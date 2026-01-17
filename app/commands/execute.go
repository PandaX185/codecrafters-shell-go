package commands

import (
	"context"
	"io"
	"os"
	"strings"
)

func ExecuteCommand(ctx context.Context, cmd string, args []string, in io.Reader, out io.Writer, errOut io.Writer) {
	switch cmd {
	case Echo.String():
		handleEcho(args, out)
		break
	case Type.String():
		cmd := strings.Join(args, " ")
		handleType(cmd, out, errOut)
		break
	case Exit.String():
		os.Exit(0)
	case Pwd.String():
		handlePwd(out, errOut)
		break
	case Cd.String():
		dir := strings.Join(args, " ")
		handleCd(dir, out, errOut)
		break
	default:
		handleExternalApp(ctx, cmd, args, in, out, errOut)
	}
	return
}
