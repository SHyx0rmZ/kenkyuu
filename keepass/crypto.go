package keepass

import (
	"bufio"
	"bytes"
	"crypto/cipher"
	"errors"
	"io"
)

type blockModeReader struct {
	BlockMode cipher.BlockMode
	Reader    io.Reader
}

func NewBlockModeReader(bm cipher.BlockMode, r io.Reader) io.ReadCloser {
	return &paddingValidator{
		Reader: bufio.NewReaderSize(&blockModeReader{
			BlockMode: bm,
			Reader:    r,
		}, bm.BlockSize()),
	}
}

func (r *blockModeReader) Read(p []byte) (n int, err error) {
	if len(p) < r.BlockMode.BlockSize() {
		return 0, io.ErrShortBuffer
	}
	d := (len(p) / r.BlockMode.BlockSize()) * r.BlockMode.BlockSize()
	n, err = r.Reader.Read(p[:d])
	r.BlockMode.CryptBlocks(p[:n], p[:n])
	return n, err
}

type blockModeWriter struct {
	BlockMode cipher.BlockMode
	Writer    io.Writer

	buf bytes.Buffer
}

func NewBlockModeWriter(bm cipher.BlockMode, w io.Writer) io.WriteCloser {
	return &blockModeWriter{
		BlockMode: bm,
		Writer:    w,
	}
}

func (w *blockModeWriter) Write(p []byte) (n int, err error) {
	n, err = w.buf.Write(p)
	if err != nil {
		return 0, err
	}
	return n, err
}

func (w *blockModeWriter) Close() error {
	b := w.BlockMode.BlockSize() - (w.buf.Len() % w.BlockMode.BlockSize())
	if b != w.BlockMode.BlockSize() {
		// Add padding as specified in PKCS#7
		w.buf.Write(bytes.Repeat([]byte{byte(b)}, b))
	}
	w.BlockMode.CryptBlocks(w.buf.Bytes(), w.buf.Bytes())
	_, err := w.Writer.Write(w.buf.Bytes())
	return err
}

type paddingValidator struct {
	*bufio.Reader
}

var ErrPadding = errors.New("invalid padding")

func (r *paddingValidator) Close() error {
	// Read the padding as defined in PKCS#7,
	// i.e. n bytes of value n.
	bs, err := r.Reader.Peek(1)
	if err != nil {
		// In case of 0 bytes, assume no
		// padding was necessary.
		if err != bufio.ErrBufferFull {
			return nil
		}
		return err
	}
	n := bs[0]
	bs = make([]byte, n)
	_, err = io.ReadFull(r.Reader, bs)
	if err != nil {
		return err
	}
	// Check all the padding bytes
	for _, b := range bs {
		if b != n {
			return ErrPadding
		}
	}
	return nil
}
