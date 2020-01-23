package inotify

import (
	"log"
	"syscall"
	"unsafe"
)
import "fmt"

type Watcher struct {
	fd int
}

func NewWatcher(flags ...Flag) (*Watcher, error) {
	var flag Flag
	for _, f := range flags {
		flag |= f
	}
	fd, err := syscall.InotifyInit1(int(flag))
	if err != nil {
		return nil, err
	}
	return &Watcher{fd}, nil
}

// Add adds a watch for the given path and mask, returning a watch
// descriptor.
func (w *Watcher) Add(path string, mask Mask) (int, error) {
	wd, err := syscall.InotifyAddWatch(w.fd, path, uint32(mask))
	if err != nil {
		return 0, err
	}

	return wd, nil
}

// Remove removes a watch that was previously registered through Add.
func (w *Watcher) Remove(wd int) error {
	_, err := syscall.InotifyRmWatch(w.fd, uint32(wd))
	return err
}

// Close closes the underlying inotify file descriptor.
func (w Watcher) Close() error {
	return syscall.Close(w.fd)
}

// Events starts a new Go-routine that sends Events related to paths
// previously registered through Add to a channel returned by this function.
func (w Watcher) Events() <-chan Event {
	ch := make(chan Event, 3)

	go func() {
		var buf [(syscall.SizeofInotifyEvent + syscall.NAME_MAX) * 3]byte
		// todo: fix race condition
		n, err := syscall.Read(w.fd, buf[:])
		if err != nil {
			log.Printf("error: %s", err)
			close(ch)
			return
		}

		for offset := 0; offset < n; {
			event := (*syscall.InotifyEvent)(unsafe.Pointer(&buf[offset]))

			ch <- fromSyscall(event)

			offset += syscall.SizeofInotifyEvent + int(event.Len)
			fmt.Println(offset, event.Len, n)
		}
	}()

	return ch
}
