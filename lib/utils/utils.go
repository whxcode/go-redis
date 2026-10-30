package utils

func ToCmdLine(args ...string) [][]byte {
	argv := make([][]byte, len(args))

	for i, arg := range args {
		argv[i] = []byte(arg)
	}

	return argv
}
