package args

import "strconv"

func Int(args []string, index int, def int) (int, error) {
	if len(args) <= index {
		return def, nil
	}

	return strconv.Atoi(args[index])
}

func String(args []string, index int, def string) string {
	if len(args) <= index {
		return def
	}

	return args[index]
}
