package ase

import (
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"
)

type ChunkType uint16

const (
	ChunkTypeOldPalette1 ChunkType = 0x0004
	ChunkTypeOldPalette2 ChunkType = 0x0011
	ChunkTypeLayer       ChunkType = 0x2004
	ChunkTypeCel         ChunkType = 0x2005
	ChunkTypeMask        ChunkType = 0x2016
	ChunkTypePath        ChunkType = 0x2017
	ChunkTypeFrameTags   ChunkType = 0x2018
	ChunkTypePalette     ChunkType = 0x2019
	ChunkTypeUserData    ChunkType = 0x2020
)

func (t ChunkType) String() string {
	switch t {
	case ChunkTypeOldPalette1:
		return "CHUNK_OLD_PALETTE1"
	case ChunkTypeOldPalette2:
		return "CHUNK_OLD_PALETTE2"
	case ChunkTypeLayer:
		return "CHUNK_LAYER"
	case ChunkTypeCel:
		return "CHUNK_CEL"
	case ChunkTypeMask:
		return "CHUNK_MASK"
	case ChunkTypePath:
		return "CHUNK_PATH"
	case ChunkTypeFrameTags:
		return "CHUNK_FRAME_TAGS"
	case ChunkTypePalette:
		return "CHUNK_PALETTE"
	case ChunkTypeUserData:
		return "CHUNK_USER_DATA"
	}

	return "CHUNK_UNKNOWN"
}

const (
	PaletteFlagName = 0x0001
)

var (
	chunkLoaderTable = map[ChunkType]func(State, io.Reader) (*ChunkBase, error){
		ChunkTypeOldPalette1: loadIgnore,
		ChunkTypeOldPalette2: loadIgnore,
		ChunkTypeLayer: func(state State, reader io.Reader) (*ChunkBase, error) {
			l, err := loadLayer(state, reader)
			return (*ChunkBase)(unsafe.Pointer(&l)), err
		},
		ChunkTypeCel: func(state State, reader io.Reader) (*ChunkBase, error) {
			c, err := loadCel(state, reader)
			return (*ChunkBase)(unsafe.Pointer(c)), err
		},
		ChunkTypeMask:      loadIgnore,
		ChunkTypePath:      loadIgnore,
		ChunkTypeFrameTags: loadIgnore,
		ChunkTypePalette: func(state State, reader io.Reader) (*ChunkBase, error) {
			p, err := loadPalette(state, reader)
			return (*ChunkBase)(unsafe.Pointer(&p)), err
		},
		ChunkTypeUserData: loadChunkUserData,
	}
)

func loadIgnore(state State, r io.Reader) (*ChunkBase, error) {
	return nil, nil
}

func loadChunkUserData(state State, reader io.Reader) (*ChunkBase, error) {
	return loadIgnore(state, reader)
}

type ChunkBase struct {
	byteSize  uint32
	chunkType ChunkType
}

func (c *ChunkBase) UnmarshalBinary(data []byte) error {
	if len(data) < 6 {
		return ErrTooShort
	}

	if c == nil {
		c = new(ChunkBase)
	}

	c.byteSize = binary.LittleEndian.Uint32(data[0:])
	c.chunkType = ChunkType(binary.LittleEndian.Uint16(data[4:]))

	return nil
}

func loadChunk(state State, r io.Reader) (c *ChunkBase, err error) {
	b := make([]byte, 6)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return
	}
	c = new(ChunkBase)
	err = c.UnmarshalBinary(b)
	if err != nil {
		return
	}
	l, ok := chunkLoaderTable[c.chunkType]
	if !ok {
		return nil, fmt.Errorf("unknown chunk type: %d", c.chunkType)
	}
	var c2 *ChunkBase
	c2, err = l(state, r)
	if err != nil {
		return
	}
	if c2 == nil {
		err = binary.Read(r, binary.LittleEndian, make([]byte, c.byteSize-6))
		return
	}
	c = c2
	err = c.UnmarshalBinary(b)
	return
}
