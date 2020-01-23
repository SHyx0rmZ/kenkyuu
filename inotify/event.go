package inotify

import (
	"syscall"
)
import "fmt"
import "unsafe"

type Event struct {
	WatchDescriptor int
	Mask            Mask
	Cookie          uint
	Name            string
}

func fromSyscall(event *syscall.InotifyEvent) Event {
	fmt.Println(event)
	bytes := make([]uint8, event.Len)
	for i := range bytes {
		bytes[i] = *(*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&event.Name))+uintptr(i)))
	}
	return Event{
		int(event.Wd),
		Mask(event.Mask),
		uint(event.Cookie),
		string(bytes),
	}
}
