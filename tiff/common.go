package tiff

import "encoding/binary"

type Bytes interface {
	Bytes(...binary.ByteOrder) []byte
}

type ByteOrder string

const (
	II ByteOrder = "II"
	MM ByteOrder = "MM"
)

// Bytes Return [2]byte
func (b ByteOrder) Bytes() []byte {
	return []byte{b[0], b[1]}
}

func (b ByteOrder) Bo() binary.ByteOrder {
	switch b {
	case II:
		return binary.LittleEndian
	case MM:
		return binary.BigEndian
	default:
		return nil
	}
}

type Version uint16

const (
	DefaultTiff Version = 0x002A
	BigTiff     Version = 0x002B
	Tiff85      Version = 0x0055
)

// Bytes Return [2]byte default ByteOrder is BigEndian
func (v Version) Bytes(bo ...binary.ByteOrder) []byte {
	var a = make([]byte, 2)
	var b binary.ByteOrder = binary.BigEndian
	if len(bo) > 0 && bo[0] == binary.LittleEndian {
		b = binary.LittleEndian
	}
	b.PutUint16(a, uint16(v))
	return a
}
