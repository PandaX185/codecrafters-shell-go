package commands

func HasRedir(args []string) (int, int) {
	for i, arg := range args {
		if arg == ">" || arg == "1>" {
			return i, 0
		}
		if arg == ">>" {
			return i, 1
		}
	}
	return -1, -1
}
