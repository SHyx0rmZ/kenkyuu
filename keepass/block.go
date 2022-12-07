package keepass

import (
	"encoding/binary"
	"io"
)

type Block struct{}

func (Block) Prefix(r io.Reader) (n int64, err error) {
	var block struct {
		Index  uint32
		Hash   [32]byte // todo: at some point, we will want to actually verify the block
		Length uint32
	}
	err = binary.Read(r, binary.LittleEndian, &block)
	if err != nil {
		return 0, err
	}
	// Check for length 0, which marks the last block.
	if block.Length == 0 {
		if c, ok := r.(io.Closer); ok {
			// We need to check for trailing padding
			// here if we're dealing with a normal
			// KeePass database.
			err = c.Close()
			if err == nil {
				err = io.EOF
			}
			return 0, err
		}
	}
	return int64(block.Length), nil
}
