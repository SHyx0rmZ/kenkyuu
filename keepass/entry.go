package keepass

import (
	"time"
)

type EntryTimes struct {
	LastModificationTime time.Time
	CreationTime         time.Time
	LastAccessTime       time.Time
	ExpiryTime           time.Time
	Expires              Bool
	UsageCount           int
	LocationChanged      time.Time
}

type StringField struct {
	Key   string
	Value StringValue
}

type StringValue struct {
	Protected Bool   `xml:",attr,omitempty"`
	Value     string `xml:",innerxml"`
}

type BinaryField struct {
	Key   string
	Value BinaryValue
}

type BinaryValue struct {
	Protected Bool   `xml:",attr,omitempty"`
	Ref       string `xml:",attr"`
	Value     string `xml:",innerxml"`
}

type AutoTypeOptions struct {
	Enabled                 Bool
	DataTransferObfuscation int
	DefaultSequence         string
}

type Entry struct {
	UUID            []byte
	IconID          int
	ForegroundColor string
	BackgroundColor string
	OverrideURL     string
	Tags            string
	Times           EntryTimes
	String          []StringField
	Binary          []BinaryField
	AutoType        AutoTypeOptions
	History         []*HistoricEntry `xml:"History>Entry"`
}

type HistoricEntry struct {
	UUID            []byte
	IconID          int
	ForegroundColor string
	BackgroundColor string
	OverrideURL     string
	Tags            string
	Times           EntryTimes
	String          []StringField
	Binary          []BinaryField
	AutoType        AutoTypeOptions
}
