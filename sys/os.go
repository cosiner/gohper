package sys

import (
	"runtime"
)

type OS string

const (
	LINUX                 OS   = "linux"
	WINDOWS                    = "windows"
	DARWIN                     = "darwin"
	FREEBSD                    = "freebsd"
	SOLARIS                    = "solaris"
	ANDROID                    = "android"
	UNKNOWN                    = "unknown"
	UNKNOWN_PATH_SEPRATOR rune = ' '
)

func (o OS) String() string {
	return string(o)
}

func OperateSystem() OS {
	return OS(runtime.GOOS)
}

func IsLinux() bool {
	return OperateSystem() == LINUX
}

func IsWindows() bool {
	return OperateSystem() == WINDOWS
}

func IsDarwin() bool {
	return OperateSystem() == DARWIN
}

func IsFreebsd() bool {
	return OperateSystem() == FREEBSD
}

func IsSolaris() bool {
	return OperateSystem() == SOLARIS
}

func IsAndroid() bool {
	return OperateSystem() == ANDROID
}

// EnvPathSeperator return seperator of env variable "PATH"
func EnvPathSeperator() rune {
	switch OperateSystem() {
	case LINUX, DARWIN, SOLARIS, FREEBSD, ANDROID:
		return ':'
	case WINDOWS:
		return ';'
	default:
		return UNKNOWN_PATH_SEPRATOR
	}
}
