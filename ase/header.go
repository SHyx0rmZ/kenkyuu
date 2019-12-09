package ase

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrTooShort = errors.New("too short")
)

type Header struct {
	fileSize         uint32
	magicNumber      uint16
	frames           uint16
	width            uint
	height           uint
	colorDepth       uint
	flags            uint32
	speed            uint16
	transparentIndex uint8
	colors           uint16
}

func (h *Header) UnmarshalBinary(p []byte) error {
	if len(p) < 128 {
		return ErrTooShort
	}

	if h == nil {
		h = new(Header)
	}

	h.fileSize = binary.LittleEndian.Uint32(p[0:])
	h.magicNumber = binary.LittleEndian.Uint16(p[4:])
	h.frames = binary.LittleEndian.Uint16(p[6:])
	h.width = uint(binary.LittleEndian.Uint16(p[8:]))
	h.height = uint(binary.LittleEndian.Uint16(p[10:]))
	h.colorDepth = uint(binary.LittleEndian.Uint16(p[12:]))
	h.flags = binary.LittleEndian.Uint32(p[14:])
	h.speed = binary.LittleEndian.Uint16(p[16:])
	h.transparentIndex = p[26]
	h.colors = binary.LittleEndian.Uint16(p[30:])

	return nil
}

func loadHeader(r io.Reader) (header Header, err error) {
	b := make([]byte, 128)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return
	}
	err = header.UnmarshalBinary(b)
	return
}
