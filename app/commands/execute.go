package commands

import "strings"

func ExecuteCommand(cmd string, args []string) (res string, err string) {
	switch cmd {
	case Echo.String():
		res, err = handleEcho(args)
		break
	case Type.String():
		cmd := strings.Join(args, " ")
		res, err = handleType(cmd)
		break
	case Exit.String():
		return
	case Pwd.String():
		res, err = handlePwd()
		break
	case Cd.String():
		dir := strings.Join(args, " ")
		res, err = handleCd(dir)
		break
	default:
		res, err = handleExternalApp(cmd, args)
	}

	return
}
