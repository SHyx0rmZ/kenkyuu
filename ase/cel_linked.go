package ase

import (
	"encoding/binary"
	"io"
)

type CelLinked struct {
	Cel
	framePosition uint16
}

func (c *CelLinked) UnmarshalBinary(data []byte) error {
	if len(data) < 2 {
		return ErrTooShort
	}

	if c == nil {
		c = new(CelLinked)
	}

	c.framePosition = binary.LittleEndian.Uint16(data[0:])

	return nil
}

func loadCelLinked(state State, r io.Reader) (c CelLinked, err error) {
	b := make([]byte, 2)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return
	}
	err = c.UnmarshalBinary(b)
	return
}
