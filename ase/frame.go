package ase

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Frame struct {
	byteSize    uint32
	magicNumber uint16
	chunks      []*ChunkBase
	duration    uint16
}

func loadFrame(state State, r io.Reader) Frame {
	var f Frame
	var l uint16

	binary.Read(r, binary.LittleEndian, &f.byteSize)
	binary.Read(r, binary.LittleEndian, &f.magicNumber)
	binary.Read(r, binary.LittleEndian, &l)
	binary.Read(r, binary.LittleEndian, &f.duration)
	binary.Read(r, binary.LittleEndian, make([]byte, 6))

	f.chunks = make([]*ChunkBase, l)

	for i := range f.chunks {
		var err error
		f.chunks[i], err = loadChunk(state, r)
		if err != nil {
			panic(fmt.Sprintf("loadFrame(%d): %s", i, err))
		}
	}

	return f
}
