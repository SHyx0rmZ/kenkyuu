package ase

import (
	"compress/zlib"
	"encoding/binary"
	"io"
)

type CelCompressed struct {
	Cel
	width  uint
	height uint
	pixels []byte
}

func (c *CelCompressed) UnmarshalBinary(data []byte) error {
	if len(data) < 4 {
		return ErrTooShort
	}

	if c == nil {
		c = new(CelCompressed)
	}

	c.width = uint(binary.LittleEndian.Uint16(data[0:]))
	c.height = uint(binary.LittleEndian.Uint16(data[2:]))

	return nil
}

func loadCelCompressed(state State, r io.Reader) (c CelCompressed, err error) {
	b := make([]byte, 4)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return
	}
	err = c.UnmarshalBinary(b)
	if err != nil {
		return
	}
	r, err = zlib.NewReader(r)
	if err != nil {
		return
	}
	defer r.(io.Closer).Close()
	c.pixels = make([]byte, c.width*c.height*state.Header.colorDepth/8)
	_, err = io.ReadFull(r, c.pixels)
	return
}
