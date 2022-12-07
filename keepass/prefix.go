package keepass

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"
)

type PrefixReader struct {
	Reader   io.Reader
	Prefixer interface {
		Prefix(r io.Reader) (n int64, err error)
	}

	r io.Reader
}

func (r *PrefixReader) Read(p []byte) (n int, err error) {
	if r.r == nil {
		n, err := r.Prefixer.Prefix(r.Reader)
		if err != nil {
			return 0, err
		}
		r.r = io.LimitReader(r.Reader, n)
	}
	n, err = r.r.Read(p)
	if err == io.EOF {
		r.r = nil
		err = nil
	}
	return n, err
}

type PrefixWriter struct {
	i int
	w io.Writer
}

func (w *PrefixWriter) Write(p []byte) (n int, err error) {
	block := struct {
		Index  uint32
		Hash   [32]byte
		Length uint32
	}{
		Index:  uint32(w.i),
		Hash:   sha256.Sum256(p),
		Length: uint32(len(p)),
	}
	w.i++
	if len(p) == 0 {
		block.Hash = [32]byte{}
	}
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, &block)
	if err != nil {
		return 0, err
	}
	n, err = buf.Write(p)
	if err != nil {
		return n, err
	}
	n2, err := io.Copy(w.w, buf)
	return int(n2), err
}
