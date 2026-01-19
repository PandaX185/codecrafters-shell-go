package commands

import (
	"os"
	"strings"
)

func ExecuteCommand(cmd string, args []string) (res string, err string) {
	switch cmd {
	case Echo.String():
		return handleEcho(args)
	case Type.String():
		cmd := strings.Join(args, " ")
		return handleType(cmd)
	case Exit.String():
		os.Exit(0)
		return
	case Pwd.String():
		return handlePwd()
	case Cd.String():
		dir := strings.Join(args, " ")
		return handleCd(dir)
	default:
		return handleExternalApp(cmd, args)
	}
}
