package ase

import (
	"encoding/binary"
	"io"
)

type PaletteEntry struct {
	flags uint16
	red   uint8
	green uint8
	blue  uint8
	alpha uint8
	name  *string
}

func (e *PaletteEntry) UnmarshalBinary(data []byte) error {
	if len(data) < 6 {
		return ErrTooShort
	}

	if e == nil {
		e = new(PaletteEntry)
	}

	e.flags = binary.LittleEndian.Uint16(data[0:])
	e.red = data[2]
	e.green = data[3]
	e.blue = data[4]
	e.alpha = data[5]

	return nil
}

func loadPaletteEntry(r io.Reader) (e PaletteEntry, err error) {
	b := make([]byte, 6)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return
	}
	err = e.UnmarshalBinary(b)
	if err != nil {
		return
	}
	if e.flags&PaletteFlagName == PaletteFlagName {
		*e.name, err = loadString(r)
		if err != nil {
			e.name = nil
		}
	}

	return
}
