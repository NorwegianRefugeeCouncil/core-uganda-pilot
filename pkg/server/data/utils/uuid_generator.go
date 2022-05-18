package utils

import (
	"github.com/nrc-no/core/pkg/server/data/api"
)

// UUIDGenerator implements the UUIDGenerator api.
// it is a zero-dependency struct that can only generate V4 uuids
type UUIDGenerator struct {
	rand api.Rand
}

// Generate implements UUIDGenerator.Generate
func (g *UUIDGenerator) Generate() (string, error) {
	var u [16]byte
	if _, err := g.rand.Read(u[:]); err != nil {
		return "", err
	}
	u[6] = (u[6] & 0x0f) | (4 << 4)
	u[8] = u[8]&(0xff>>2) | (0x02 << 6)
	buf := make([]byte, 36)
	encodeHex(buf[0:8], u[0:4])
	buf[8] = '-'
	encodeHex(buf[9:13], u[4:6])
	buf[13] = '-'
	encodeHex(buf[14:18], u[6:8])
	buf[18] = '-'
	encodeHex(buf[19:23], u[8:10])
	buf[23] = '-'
	encodeHex(buf[24:], u[10:])
	return string(buf), nil
}

const hexTable = "0123456789abcdef"

// EncodeHex encodes a byte array to hexadecimal string
// dst is the destination buffer
// src is the source buffer
// it returns the number of bytes encoded in dst
// it is a zero-dependent version of hex.Encode
func encodeHex(dst, src []byte) int {
	j := 0
	for _, v := range src {
		dst[j] = hexTable[v>>4]
		dst[j+1] = hexTable[v&0x0f]
		j += 2
	}
	return len(src) * 2
}
