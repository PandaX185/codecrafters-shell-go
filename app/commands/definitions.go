package commands

type command int

const (
	Echo command = iota
	Exit
	Type
)

var builtinCommands = map[string]bool{
	Echo.String():  true,
	Exit.String():  true,
	Type.String(): true,
}

func (c command) String() string {
	switch c {
	case Echo:
		return "echo"
	case Exit:
		return "exit"
	case Type:
		return "type"
	default:
		return "unknown"
	}
}
