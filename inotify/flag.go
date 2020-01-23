package inotify

import (
	"syscall"
)

type Flag int

const (
	FlagCloseOnExec Flag = syscall.IN_CLOEXEC
	FlagNonBlocking Flag = syscall.IN_NONBLOCK
)
