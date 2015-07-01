package args

import "strconv"

// Int get arguments at given index, if not exists, use default value,
//  otherwise, convert it to inteeger
func Int(args []string, index int, def int) (int, error) {
	if len(args) <= index {
		return def, nil
	}

	return strconv.Atoi(args[index])
}

// String get argument at given index, if not exists, use default value
func String(args []string, index int, def string) string {
	if len(args) <= index {
		return def
	}

	return args[index]
}
