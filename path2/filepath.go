package path2

import "os"

func RemoveExt(path string) string {
	for i := len(path) - 1; i >= 0 && !os.IsPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			return path[:i]
		}
	}
	return path
}

func ReplaceExt(path, ext string) string {
	return RemoveExt(path) + ext
}
