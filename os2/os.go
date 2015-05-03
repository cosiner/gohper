package os2

import (
	"runtime"
)

const (
	LINUX   = "linux"
	WINDOWS = "windows"
	DARWIN  = "darwin"
	FREEBSD = "freebsd"
	SOLARIS = "solaris"
	ANDROID = "android"
	UNKNOWN = "unknown"
)

func OS() string {
	return runtime.GOOS
}

func IsLinux() bool {
	return OS() == LINUX
}

func IsWindows() bool {
	return OS() == WINDOWS
}

func IsDarwin() bool {
	return OS() == DARWIN
}

func IsFreebsd() bool {
	return OS() == FREEBSD
}

func IsSolaris() bool {
	return OS() == SOLARIS
}

func IsAndroid() bool {
	return OS() == ANDROID
}
