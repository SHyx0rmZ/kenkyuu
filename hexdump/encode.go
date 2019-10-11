package hexdump

import (
	"encoding/hex"
	"io"
)

func NewEncoder(w io.Writer) io.WriteCloser {
	return hex.Dumper(w)
}
