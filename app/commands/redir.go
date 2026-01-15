package commands

func HasOutRedir(args []string) (int, int) {
	for i, arg := range args {
		if arg == ">" || arg == "1>" {
			return i, 0
		}
		if arg == ">>" || arg == "1>>" {
			return i, 1
		}
	}
	return -1, -1
}

func HasErrRedir(args []string) (int, int) {
	for i, arg := range args {
		if arg == "2>" {
			return i, 0
		}
		if arg == "2>>" {
			return i, 1
		}
	}
	return -1, -1
}
