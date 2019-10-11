package hexdump

import (
	"encoding/hex"
	"io"
)

const bytesPerLine = 16
const charsPerByte = 3

type decoder struct {
	r   io.Reader
	buf [79]byte
}

func NewDecoder(r io.Reader) io.Reader {
	return &decoder{r, [79]byte{}}
}

func (d *decoder) Read(p []byte) (n int, err error) {
	// Read a single line into a buffer, then decode its
	// bytes into the beginning of the buffer, ignoring the
	// offset and ASCII representation.
	n, err = io.ReadAtLeast(d.r, d.buf[:], 62)
	if err != nil {
		return
	}
	// There are 16 bytes per line, each represented by two
	// characters and delimited by a single space. There is
	// an additional space after the eighth byte. The bytes
	// start at offset 10 in each line.
	for i := 0; i < bytesPerLine; i++ {
		off := i*charsPerByte + (i / (bytesPerLine / 2)) + 10
		if d.buf[off] == ' ' {
			n = copy(p, d.buf[:i])
			return n, nil
		}
		n, err = hex.Decode(d.buf[i:], d.buf[off:off+2])
		if err != nil {
			return
		}
	}
	// todo: We need to check whether we're dealing with malformed
	//       input here. There might be more unexpected bytes.
	n = copy(p, d.buf[:bytesPerLine])
	return n, nil
}
