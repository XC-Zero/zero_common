package tiff

import (
	"encoding/binary"
	"math"
	"slices"
)

type String string

// Bytes Return [2]byte default ByteOrder is BigEndian
func (s String) Bytes(bo ...binary.ByteOrder) []byte {
	var b = []byte(string(s))
	if len(bo) > 0 && bo[0] == binary.LittleEndian {
		slices.Reverse(b)
	}

	return b
}

type Byte byte

func (b Byte) Bytes(...binary.ByteOrder) []byte {
	return []byte{byte(b)}
}

type Ascii byte

func (b Ascii) Bytes(...binary.ByteOrder) []byte {
	return []byte{byte(b)}
}

// Short SHORT
type Short uint16

// Bytes Return [2]byte default ByteOrder is BigEndian
func (s Short) Bytes(bo ...binary.ByteOrder) []byte {
	var n = make([]byte, 2)
	var b binary.ByteOrder = binary.BigEndian
	if len(bo) > 0 && bo[0] == binary.LittleEndian {
		b = binary.LittleEndian
	}
	b.PutUint16(n, uint16(s))
	return n
}

// Long LONG
type Long uint32

// Bytes Return [4]byte default ByteOrder is BigEndian
func (l Long) Bytes(bo ...binary.ByteOrder) []byte {
	var n = make([]byte, 4)
	var b binary.ByteOrder = binary.BigEndian
	if len(bo) > 0 && bo[0] == binary.LittleEndian {
		b = binary.LittleEndian
	}
	b.PutUint32(n, uint32(l))
	return n
}

// Rational RATIONAL Two LONGs: the first represents the numerator
type Rational struct {
	numerator uint32
	fraction  uint32
}

// Bytes Return [8]byte default ByteOrder is BigEndian
func (r Rational) Bytes(bo ...binary.ByteOrder) []byte {
	var n, f = make([]byte, 4), make([]byte, 4)
	var b binary.ByteOrder = binary.BigEndian
	if len(bo) > 0 && bo[0] == binary.LittleEndian {
		b = binary.LittleEndian
	}
	b.PutUint32(n, r.numerator)
	b.PutUint32(f, r.fraction)

	return append(n, f...)
}

type Float float32

// Bytes Return [4]byte default ByteOrder is BigEndian
func (f Float) Bytes(bo ...binary.ByteOrder) []byte {
	var n = make([]byte, 4)
	var b binary.ByteOrder = binary.BigEndian
	if len(bo) > 0 && bo[0] == binary.LittleEndian {
		b = binary.LittleEndian
	}

	b.PutUint32(n, math.Float32bits(float32(f)))
	return n
}

type Double float32

// Bytes Return [8]byte default ByteOrder is BigEndian
func (d Double) Bytes(bo ...binary.ByteOrder) []byte {
	var n = make([]byte, 8)
	var b binary.ByteOrder = binary.BigEndian
	if len(bo) > 0 && bo[0] == binary.LittleEndian {
		b = binary.LittleEndian
	}

	b.PutUint64(n, math.Float64bits(float64(d)))
	return n
}

type ComplementFour struct {
	B Bytes
}

func (c ComplementFour) Bytes(bo ...binary.ByteOrder) []byte {
	var n = make([]byte, 0, 4)

	var b binary.ByteOrder = binary.BigEndian
	if len(bo) > 0 && bo[0] == binary.LittleEndian {
		b = binary.LittleEndian
	}

	var res = c.B.Bytes(b)
	if len(res) == 4 {
		n = res
	}
	if len(res) < 4 {
		n = res
		for len(n) < 4 {
			n = append(n, 0)
		}
	}
	if len(res) > 4 {
		return res[:4]
	}
	return n
}
