package commands

import "strings"

func ExecuteCommand(cmd string, args []string) (res string, err string) {
	switch cmd {
	case Echo.String():
		return handleEcho(args)
	case Type.String():
		cmd := strings.Join(args, " ")
		return handleType(cmd)
	case Exit.String():
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
