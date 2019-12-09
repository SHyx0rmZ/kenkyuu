package ase

import (
	"encoding/binary"
	"io"
)

type Palette struct {
	ChunkBase
	paletteSize uint32
	firstColor  uint32
	lastColor   uint32
	entries     []PaletteEntry
}

func (p *Palette) UnmarshalBinary(data []byte) error {
	if len(data) < 20 {
		return ErrTooShort
	}

	if p == nil {
		p = new(Palette)
	}

	p.paletteSize = binary.LittleEndian.Uint32(data[0:])
	p.firstColor = binary.LittleEndian.Uint32(data[4:])
	p.lastColor = binary.LittleEndian.Uint32(data[8:])

	return nil
}

func loadPalette(state State, r io.Reader) (p Palette, err error) {
	b := make([]byte, 20)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return
	}
	err = p.UnmarshalBinary(b)
	if err != nil {
		return
	}
	for i := p.lastColor - p.firstColor + 1; i > 0; i-- {
		var e PaletteEntry
		e, err = loadPaletteEntry(r)
		if err != nil {
			return
		}
		p.entries = append(p.entries, e)
	}
	return
}
