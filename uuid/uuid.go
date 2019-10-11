package uuid

import (
	"math/rand"
)

const hextable = "0123456789abcdef"

type UUID [16]byte

func New() UUID {
	var u UUID
	rand.Read(u[:])
	u[6] = 0x40 | u[6]>>4
	u[8] |= 0x80
	return u
}

func (bs UUID) String() string {
	return string([]byte{
		hextable[bs[0]>>4],
		hextable[bs[0]&15],
		hextable[bs[1]>>4],
		hextable[bs[1]&15],
		hextable[bs[2]>>4],
		hextable[bs[2]&15],
		hextable[bs[3]>>4],
		hextable[bs[3]&15],
		'-',
		hextable[bs[4]>>4],
		hextable[bs[4]&15],
		hextable[bs[5]>>4],
		hextable[bs[5]&15],
		'-',
		hextable[bs[6]>>4],
		hextable[bs[6]&15],
		hextable[bs[7]>>4],
		hextable[bs[7]&15],
		'-',
		hextable[bs[8]>>4],
		hextable[bs[8]&15],
		hextable[bs[9]>>4],
		hextable[bs[9]&15],
		'-',
		hextable[bs[10]>>4],
		hextable[bs[10]&15],
		hextable[bs[11]>>4],
		hextable[bs[11]&15],
		hextable[bs[12]>>4],
		hextable[bs[12]&15],
		hextable[bs[13]>>4],
		hextable[bs[13]&15],
		hextable[bs[14]>>4],
		hextable[bs[14]&15],
		hextable[bs[15]>>4],
		hextable[bs[15]&15],
	})
}
