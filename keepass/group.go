package keepass

import (
	"time"
)

type Group struct {
	UUID   []byte
	Name   string
	Notes  string
	IconID int
	Times  struct {
		LastModificationTime time.Time
		CreationTime         time.Time
		LastAccessTime       time.Time
		ExpiryTime           time.Time
		Expires              Bool
		UsageCount           int
		LocationChanged      time.Time
	}
	IsExpanded              Bool
	DefaultAutoTypeSequence string
	EnableAutoType          string
	EnableSearching         string
	LastTopVisibleEntry     string
	Group                   []*Group
	Entry                   []*Entry
}
