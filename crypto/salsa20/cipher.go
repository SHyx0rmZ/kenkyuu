package salsa20

import (
	"bytes"
	"crypto/cipher"
	"unsafe"

	"golang.org/x/crypto/salsa20/salsa"
)

func New(nonce *[8]byte, key *[32]byte) cipher.Stream {
	s := &stream{
		counter: &[16]byte{},
		key:     key,
		buf:     &[BlockSize]byte{},
	}
	copy(s.counter[:], nonce[:])
	return s
}

type stream struct {
	counter *[NonceSize]byte
	key     *[KeySize]byte
	buf     *[BlockSize]byte
	pos     int
}

const (
	NonceSize     = 16
	KeySize       = 32
	BlockSize     = 64
	CounterOffset = 8
)

var emptyBuffer = bytes.Repeat([]byte{0}, BlockSize)

func (c *stream) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic("salsa20: output smaller than input")
	}
	if InexactOverlap(dst[:len(src)], src) {
		panic("salsa20: invalid buffer overlap")
	}
	block := c.pos / BlockSize
	c.updateCounter(block)
	off := c.pos % 64
	if off == 0 {
		salsa.XORKeyStream(dst, src, c.counter, c.key)
		if (c.pos + len(src)) < c.pos {
			panic("overflow")
		}
		c.pos += len(src)
		return
	}
	n := copy(c.buf[off:], src)
	salsa.XORKeyStream(c.buf[:off+n], c.buf[:off+n], c.counter, c.key)
	copy(dst, c.buf[off:off+n])
	copy(c.buf[:], emptyBuffer)
	c.updateCounter(block + 1)
	salsa.XORKeyStream(dst[n:], src[n:], c.counter, c.key)
	if (c.pos + len(src)) < c.pos {
		panic("overflow")
	}
	c.pos += len(src)
}

func (c *stream) updateCounter(val int) {
	// put little-endian counter into nonce
	//for i := CounterOffset; i < NonceSize; i++ {
	//	c.counter[i] = byte(val)
	//	val >>= 8
	//}
	*(*uint64)(unsafe.Pointer(&c.counter[CounterOffset])) = uint64(val)
}

// AnyOverlap reports whether x and y share memory at any (not necessarily
// corresponding) index. The memory beyond the slice length is ignored.
func AnyOverlap(x, y []byte) bool {
	return len(x) > 0 && len(y) > 0 &&
		uintptr(unsafe.Pointer(&x[0])) <= uintptr(unsafe.Pointer(&y[len(y)-1])) &&
		uintptr(unsafe.Pointer(&y[0])) <= uintptr(unsafe.Pointer(&x[len(x)-1]))
}

// InexactOverlap reports whether x and y share memory at any non-corresponding
// index. The memory beyond the slice length is ignored. Note that x and y can
// have different lengths and still not have any inexact overlap.
//
// InexactOverlap can be used to implement the requirements of the crypto/cipher
// AEAD, Block, BlockMode and Stream interfaces.
func InexactOverlap(x, y []byte) bool {
	if len(x) == 0 || len(y) == 0 || &x[0] == &y[0] {
		return false
	}
	return AnyOverlap(x, y)
}
