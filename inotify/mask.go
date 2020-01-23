package inotify

import "syscall"

type Mask uint32

const (
	MaskAccessed         Mask = syscall.IN_ACCESS
	MaskAttributeChanged Mask = syscall.IN_ATTRIB
	MaskWriteClosed      Mask = syscall.IN_CLOSE_WRITE
	MaskNoWriteClosed    Mask = syscall.IN_CLOSE_NOWRITE
	MaskCreated          Mask = syscall.IN_CREATE
	MaskDeletedFrom      Mask = syscall.IN_DELETE
	MaskDeleted          Mask = syscall.IN_DELETE_SELF
	MaskModified         Mask = syscall.IN_MODIFY
	MaskMoved            Mask = syscall.IN_MOVE_SELF
	MaskMovedFrom        Mask = syscall.IN_MOVED_FROM
	MaskMovedTo          Mask = syscall.IN_MOVED_TO
	MaskOpened           Mask = syscall.IN_OPEN

	MaskDontFollow        Mask = syscall.IN_DONT_FOLLOW
	MaskExcludingUnlinked Mask = syscall.IN_EXCL_UNLINK
	MaskOneshot           Mask = syscall.IN_ONESHOT
	MaskDirectoryOnly     Mask = syscall.IN_ONLYDIR

	MaskIgnored              Mask = syscall.IN_IGNORED
	MaskIsDirectory          Mask = syscall.IN_ISDIR
	MaskEventQueueOverflowed Mask = syscall.IN_Q_OVERFLOW
	MaskUnmounted            Mask = syscall.IN_UNMOUNT
)
