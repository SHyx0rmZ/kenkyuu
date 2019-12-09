package ase

import (
	"encoding/binary"
	"io"
)

func loadString(r io.Reader) (string, error) {
	b := make([]byte, 2)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return "", err
	}
	l := binary.LittleEndian.Uint16(b)
	b = make([]byte, l)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
