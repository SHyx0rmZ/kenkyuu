package ase

import (
	"encoding/binary"
	"io"
)

type Layer struct {
	ChunkBase
	flags         uint16
	layerType     uint16
	childLevel    uint16
	defaultWidth  uint16
	defaultHeight uint16
	blendMode     uint16
	opacity       uint8
	name          string
}

func (l *Layer) UnmarshalBinary(data []byte) error {
	if len(data) < 16 {
		return ErrTooShort
	}

	if l == nil {
		l = new(Layer)
	}

	l.flags = binary.LittleEndian.Uint16(data[0:])
	l.layerType = binary.LittleEndian.Uint16(data[2:])
	l.childLevel = binary.LittleEndian.Uint16(data[4:])
	l.defaultWidth = binary.LittleEndian.Uint16(data[6:])
	l.defaultHeight = binary.LittleEndian.Uint16(data[8:])
	l.blendMode = binary.LittleEndian.Uint16(data[10:])
	l.opacity = data[12]

	return nil
}

func loadLayer(state State, r io.Reader) (l Layer, err error) {
	b := make([]byte, 16)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return
	}
	err = l.UnmarshalBinary(b)
	if err != nil {
		return
	}
	l.name, err = loadString(r)
	return
}
