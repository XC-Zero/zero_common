package tiff

import (
	"encoding/binary"
)

// DataType It is defined by TIFF6.0
type DataType uint16

// Bytes Return [2]byte default ByteOrder is BigEndian
func (d DataType) Bytes(bo ...binary.ByteOrder) []byte {
	var a = make([]byte, 2)
	var b binary.ByteOrder = binary.BigEndian
	if len(bo) > 0 && bo[0] == binary.LittleEndian {
		b = binary.LittleEndian
	}
	b.PutUint16(a, uint16(d))
	return a[:]
}

const (
	//	BYTE 8-bit unsigned integer
	//		In Go type is byte
	BYTE DataType = iota + 1
	//	ASCII 8-bit byte that contains a 7-bit ASCII code; the last byte must be NUL (binary zero).
	ASCII
	//	SHORT 16-bit (2-byte) unsigned integer.
	//		In Go type is uint16
	SHORT
	//	LONG 32-bit (4-byte) unsigned integer.
	//		In Go type is uint32
	LONG
	//	RATIONAL Two LONGs: the first represents the numerator
	RATIONAL
	//	SBYTE An 8-bit signed (twos-complement) integer.
	//		In Go type is int8
	SBYTE
	//	UNDEFINED An 8-bit byte that may contain anything, depending on the definition of the field.
	UNDEFINED
	//	SSHORT A 16-bit (2-byte) signed (twos-complement) integer.
	//		In Go type is int16
	SSHORT
	//	SLONG A 32-bit (4-byte) signed (twos-complement) integer.
	//		In Go type is int32
	SLONG
	//	SRATIONAL Two SLONG’s: the first represents the numerator of a fraction, the second the denominator.
	SRATIONAL
	//	FLOAT Single precision (4-byte) IEEE format.
	//		In Go type is float32
	FLOAT
	//	DOUBLE double precision (8-byte) IEEE format.
	//		In Go type is float64
	DOUBLE
)

func (d DataType) Len() uint32 {
	switch d {
	case BYTE, ASCII, SBYTE, UNDEFINED:
		return 1
	case SHORT, SSHORT:
		return 2
	case LONG, SLONG, FLOAT:
		return 4
	case RATIONAL, SRATIONAL, DOUBLE:
		return 8
	}
	return 0
}

type GeoType DataType

const (
	LinearMeter                     GeoType = 9001
	LinearYardIndian                GeoType = 9013
	LinearFoot                      GeoType = 9002
	LinearFathom                    GeoType = 9014
	LinearFootUSSurvey              GeoType = 9003
	LinearMileInternationalNautical GeoType = 9015
	LinearFootModifiedAmerican      GeoType = 9004
	AngularRadian                   GeoType = 9101
	LinearFootClarke                GeoType = 9005
	AngularDegree                   GeoType = 9102
	LinearFootIndian                GeoType = 9006
	AngularArcMinute                GeoType = 9103
	LinearLink                      GeoType = 9007
	AngularArcSecond                GeoType = 9104
	LinearLinkBenoit                GeoType = 9008
	AngularGrad                     GeoType = 9105
	LinearLinkSears                 GeoType = 9009
	AngularGon                      GeoType = 9106
	LinearChainBenoit               GeoType = 9010
	AngularDMSGeoType               GeoType = 9107
	LinearChainSears                GeoType = 9011
	AngularDMSHemisphere            GeoType = 9108
	LinearYardSears                 GeoType = 9012
	UserDefinedGeoType              GeoType = 32767
)
