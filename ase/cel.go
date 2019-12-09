package ase

import (
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"
)

type CelType uint16

const (
	CelTypeRaw CelType = iota
	CelTypeLinked
	CelTypeCompressed
)

var (
	celLoaderTable = map[CelType]func(State, io.Reader) (*Cel, error){
		CelTypeRaw: func(state State, r io.Reader) (*Cel, error) {
			c, err := loadCelRaw(state, r)
			return (*Cel)(unsafe.Pointer(&c)), err
		},
		CelTypeLinked: func(state State, r io.Reader) (*Cel, error) {
			c, err := loadCelLinked(state, r)
			return (*Cel)(unsafe.Pointer(&c)), err
		},
		CelTypeCompressed: func(state State, r io.Reader) (*Cel, error) {
			c, err := loadCelCompressed(state, r)
			return (*Cel)(unsafe.Pointer(&c)), err
		},
	}
)

type Cel struct {
	ChunkBase
	layerIndex uint16
	positionX  uint16
	positionY  uint16
	opacity    uint8
	celType    CelType
}

func (c *Cel) UnmarshalBinary(data []byte) error {
	if len(data) < 16 {
		return ErrTooShort
	}

	if c == nil {
		c = new(Cel)
	}

	c.layerIndex = binary.LittleEndian.Uint16(data[0:])
	c.positionX = binary.LittleEndian.Uint16(data[2:])
	c.positionY = binary.LittleEndian.Uint16(data[4:])
	c.opacity = data[6]
	c.celType = CelType(binary.LittleEndian.Uint16(data[7:]))

	return nil
}

func loadCel(state State, r io.Reader) (c *Cel, err error) {
	b := make([]byte, 16)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return
	}
	c = new(Cel)
	err = c.UnmarshalBinary(b)
	if err != nil {
		return
	}
	l, ok := celLoaderTable[c.celType]
	if !ok {
		return nil, fmt.Errorf("unknown cel type: %d", c.celType)
	}
	c, err = l(state, r)
	if err != nil {
		return
	}
	err = c.UnmarshalBinary(b)
	return
}
