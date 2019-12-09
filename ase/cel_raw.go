package ase

import (
	"encoding/binary"
	"io"
)

type CelRaw struct {
	Cel
	width  uint
	height uint
	pixels []byte
}

func (c *CelRaw) UnmarshalBinary(data []byte) error {
	if len(data) < 4 {
		return ErrTooShort
	}

	if c == nil {
		c = new(CelRaw)
	}

	c.width = uint(binary.LittleEndian.Uint16(data[0:]))
	c.height = uint(binary.LittleEndian.Uint16(data[2:]))

	return nil
}

func loadCelRaw(state State, r io.Reader) (c CelRaw, err error) {
	b := make([]byte, 4)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return
	}
	err = c.UnmarshalBinary(b)
	if err != nil {
		return
	}
	c.pixels = make([]byte, c.width*c.height*state.Header.colorDepth/8)
	_, err = io.ReadFull(r, b)
	return
}
