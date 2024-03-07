package tiff

import (
	"bytes"
	"compress/lzw"
	"encoding/binary"
	"image"
	"image/jpeg"
	"io"
	"sync"

	"github.com/pkg/errors"
)

// Tag Doc in https://www.awaresystems.be/imaging/tiff.html
type Tag uint16

// Bytes Return [2]byte default ByteOrder is BigEndian
func (t Tag) Bytes(bo ...binary.ByteOrder) []byte {
	var a = make([]byte, 2)
	var b binary.ByteOrder = binary.BigEndian
	if len(bo) > 0 && bo[0] == binary.LittleEndian {
		b = binary.LittleEndian
	}
	b.PutUint16(a, uint16(t))
	return a[:]
}

const (
	// NextSubFileType LONG
	NextSubFileType Tag = 0x00FE
	// SubFileType SHORT
	SubFileType Tag = 0x00FF
	// ImageWidth LONG
	ImageWidth Tag = 0x0100
	// ImageLength LONG
	ImageLength Tag = 0x0101
	// BitsPerSample SHORT
	BitsPerSample Tag = 0x0102
	// Compression SHORT
	Compression Tag = 0x0103
	// PhotoMetricInterpretation SHORT
	PhotoMetricInterpretation Tag = 0x0106
	// Threshholding SHORT
	Threshholding Tag = 0x0107
	CellWidth     Tag = 0x0108
	CellLength    Tag = 0x0109
	// FillOrder SHORT
	FillOrder Tag = 0x010A
	// ImageDescription ASCII
	ImageDescription Tag = 0x010E
	// Make ASCII
	Make Tag = 0x010F
	// Model ASCII
	Model Tag = 0x0110
	// StripOffsets LONG
	StripOffsets Tag = 0x0111
	// Orient SHORT
	Orient Tag = 0x0112
	// SamplesPerPixel SHORT
	SamplesPerPixel Tag = 0x0115
	// RowsPerStrip LONG
	RowsPerStrip Tag = 0x0116
	// StripByteCounts LONG
	StripByteCounts Tag = 0x0117
	// MinSampleValue SHORT
	MinSampleValue Tag = 0x0118
	// MaxSampleValue SHORT
	MaxSampleValue Tag = 0x0119
	// XResolution RATIONAL
	XResolution Tag = 0x011A
	// YResolution RATIONAL
	YResolution Tag = 0x011B
	// PlanarConfiguration SHORT
	PlanarConfiguration Tag = 0x011C
	// FreeOffset LONG
	FreeOffset Tag = 0x0120
	// FreeByteCounts LONG
	FreeByteCounts Tag = 0x0121
	// GrayResponseUnit SHORT
	GrayResponseUnit Tag = 0x0122
	// GrayResponseCurve SHORT
	GrayResponseCurve Tag = 0x0123
	// ResolutionUnit   SHORT
	ResolutionUnit Tag = 0x0128
	//Software  ASCII
	Software Tag = 0x0131
	// Datetime ASCII
	Datetime Tag = 0x0132
	// Artist ASCII
	Artist Tag = 0x013B
	// HostComputer ASCII
	HostComputer Tag = 0x013C
	// ColorMap SHORT RGB
	ColorMap      Tag = 0x0140
	TileWidth     Tag = 0x0142
	TileLength    Tag = 0x0143
	TileOffset    Tag = 0x0144
	TileByteCount Tag = 0x0145
	SubIFD        Tag = 0x014a
	// ExtraSamples SHORT
	ExtraSamples        Tag = 0x0152
	JPEGTables          Tag = 0x0015b
	GlobalParametersIFD Tag = 0x00190
	YCbCrCoefficients   Tag = 0x00211
	YCbCrSubSampling    Tag = 0x00212
	YCbCrPositioning    Tag = 0x00213
	// Copyright ASCII
	Copyright Tag = 0x8298
)

var tagNameMap = map[Tag]string{
	NextSubFileType:           "NextSubFileType",
	SubFileType:               "SubFileType",
	ImageWidth:                "ImageWidth",
	ImageLength:               "ImageLength",
	BitsPerSample:             "BitsPerSample",
	Compression:               "Compression",
	PhotoMetricInterpretation: "PhotoMetricInterpretation",
	Threshholding:             "Threshholding",
	CellWidth:                 "CellWidth",
	CellLength:                "CellLength",
	FillOrder:                 "FillOrder",
	ImageDescription:          "ImageDescription",
	Make:                      "Make",
	Model:                     "Model",
	StripOffsets:              "StripOffsets",
	Orient:                    "Orient",
	SamplesPerPixel:           "SamplesPerPixel",
	RowsPerStrip:              "RowsPerStrip",
	StripByteCounts:           "StripByteCounts",
	MinSampleValue:            "MinSampleValue",
	MaxSampleValue:            "MaxSampleValue",
	XResolution:               "XResolution",
	YResolution:               "YResolution",
	PlanarConfiguration:       "PlanarConfiguration",
	FreeOffset:                "FreeOffset",
	FreeByteCounts:            "FreeByteCounts",
	GrayResponseUnit:          "GrayResponseUnit",
	GrayResponseCurve:         "GrayResponseCurve",
	ResolutionUnit:            "ResolutionUnit",
	Software:                  "Software",
	Datetime:                  "Datetime",
	Artist:                    "Artist",
	HostComputer:              "HostComputer",
	ColorMap:                  "ColorMap",
	TileWidth:                 "TileWidth",
	TileLength:                "TileLength",
	TileOffset:                "TileOffset",
	TileByteCount:             "TileByteCount",
	SubIFD:                    "SubIFD",
	ExtraSamples:              "ExtraSamples",
	JPEGTables:                "JPEGTables",
	GlobalParametersIFD:       "GlobalParametersIFD",
	YCbCrCoefficients:         "YCbCrCoefficients",
	YCbCrSubSampling:          "YCbCrSubSampling",
	YCbCrPositioning:          "YCbCrPositioning",
	Copyright:                 "Copyright",
}

func (t Tag) String() string {
	if v, ok := tagNameMap[t]; ok {
		return v
	}
	return "UNKNOWN"
}

// GeoTag https://exiftool.org/TagNames/GeoTiff.html
type GeoTag Tag

const (
	GeoTiffVersion           GeoTag = 0x0001
	GTModelType              GeoTag = 0x0400
	GTRasterType             GeoTag = 0x0401
	GTCitation               GeoTag = 0x0402
	GeographicType           GeoTag = 0x0800
	GeogCitation             GeoTag = 0x0801
	GeogGeodeticDatum        GeoTag = 0x0802
	GeogPrimeMeridian        GeoTag = 0x0803
	GeogLinearUnits          GeoTag = 0x0804
	GeogLinearUnitSize       GeoTag = 0x0805
	GeogAngularUnits         GeoTag = 0x0806
	GeogAngularUnitSize      GeoTag = 0x0807
	GeogEllipsoid            GeoTag = 0x0808
	GeogSemiMajorAxis        GeoTag = 0x0809
	GeogSemiMinorAxis        GeoTag = 0x080a
	GeogInvFlattening        GeoTag = 0x080b
	GeogAzimuthUnits         GeoTag = 0x080c
	GeogPrimeMeridianLong    GeoTag = 0x080d
	GeogToWGS84              GeoTag = 0x080e
	ProjectedCSType          GeoTag = 0x0c00
	PCSCitation              GeoTag = 0x0c01
	Projection               GeoTag = 0x0c02
	ProjCoordTrans           GeoTag = 0x0c03
	ProjLinearUnits          GeoTag = 0x0c04
	ProjLinearUnitSize       GeoTag = 0x0c05
	ProjStdParallel1         GeoTag = 0x0c06
	ProjStdParallel2         GeoTag = 0x0c07
	ProjNatOriginLong        GeoTag = 0x0c08
	ProjNatOriginLat         GeoTag = 0x0c09
	ProjFalseEasting         GeoTag = 0x0c0a
	ProjFalseNorthing        GeoTag = 0x0c0b
	ProjFalseOriginLong      GeoTag = 0x0c0c
	ProjFalseOriginLat       GeoTag = 0x0c0d
	ProjFalseOriginEasting   GeoTag = 0x0c0e
	ProjFalseOriginNorthing  GeoTag = 0x0c0f
	ProjCenterLong           GeoTag = 0x0c10
	ProjCenterLat            GeoTag = 0x0c11
	ProjCenterEasting        GeoTag = 0x0c12
	ProjCenterNorthing       GeoTag = 0x0c13
	ProjScaleAtNatOrigin     GeoTag = 0x0c14
	ProjScaleAtCenter        GeoTag = 0x0c15
	ProjAzimuthAngle         GeoTag = 0x0c16
	ProjStraightVertPoleLong GeoTag = 0x0c17
	ProjRectifiedGridAngle   GeoTag = 0x0c18
	VerticalCSType           GeoTag = 0x1000
	VerticalCitation         GeoTag = 0x1001
	VerticalDatum            GeoTag = 0x1002
	VerticalUnits            GeoTag = 0x1003

	ModelPixelScaleTag     GeoTag = 0x0830E
	IPTC                   GeoTag = 0x083BB
	INGRPacketDataTag      GeoTag = 0x0847E
	INGRFlagRegisters      GeoTag = 0x0847F
	InterGraphMatrixTag    GeoTag = 0x08480
	ModelTiePointTag       GeoTag = 0x08482
	ModelTransformationTag GeoTag = 0x085D8
	Photoshop              GeoTag = 0x08649
	ExifIFD                GeoTag = 0x08769
	ICCProfile             GeoTag = 0x08773
	GeoKeyDirectoryTag     GeoTag = 0x087AF
	GeoDoubleParamsTag     GeoTag = 0x087B0
	GeoAsciiParamsTag      GeoTag = 0x087B1
	GPSIFD                 GeoTag = 0x08825
	HylaFAXFaxRecvParams   GeoTag = 0x0885C
	HylaFAXFaxSubAddress   GeoTag = 0x0885D
	HylaFAXFaxRecvTime     GeoTag = 0x0885E
	ImageSourceData        GeoTag = 0x0935C
	InteroperabilityIFD    GeoTag = 0x0A005
	GDAL_METADATA          GeoTag = 0xA480
	GDAL_NODATA            GeoTag = 0xA481

	ChartFormat          GeoTag = 0xb799
	ChartSource          GeoTag = 0xb79a
	ChartSourceEdition   GeoTag = 0xb79b
	ChartSourceDate      GeoTag = 0xb79c
	ChartCorrDate        GeoTag = 0xb79d
	ChartCountryOrigin   GeoTag = 0xb79e
	ChartRasterEdition   GeoTag = 0xb79f
	ChartSoundingDatum   GeoTag = 0xb7a0
	ChartDepthUnits      GeoTag = 0xb7a1
	ChartMagVar          GeoTag = 0xb7a2
	ChartMagVarYear      GeoTag = 0xb7a3
	ChartMagVarAnnChange GeoTag = 0xb7a4
	ChartWGSNSShift      GeoTag = 0xb7a5
	InsetNWPixelX        GeoTag = 0xb7a7
	InsetNWPixelY        GeoTag = 0xb7a8
	ChartContourInterval GeoTag = 0xb7a9
)

var (
	TypeErr    = errors.New("Tag value type is incorrect!")
	InvalidErr = errors.New("Tag value is invalid!")
)

type Validator interface {
	Valid(any any) error
}

type TagInfo struct {
	Tag      Tag
	DataType DataType
	Validator
}

var TagInfoMap = map[Tag]TagInfo{
	Compression: {
		Tag:       Compression,
		DataType:  SHORT,
		Validator: &defaultCompression,
	},
}

type defaultValidator struct {
	sync.Once
	vv map[TagValue]struct{}
	v  []TagValue
}

func (c *defaultValidator) Valid(a any) error {
	c.Do(func() {
		c.vv = make(map[TagValue]struct{})
		for _, value := range c.v {
			c.vv[value] = struct{}{}
		}
	})
	v, ok := a.(uint16)
	if !ok {
		return errors.WithStack(TypeErr)
	}
	if _, ok := c.vv[TagValue(v)]; !ok {
		return errors.WithStack(InvalidErr)
	}
	return nil
}

type TagValue uint16

// Bytes Return [2]byte default ByteOrder is BigEndian
func (t TagValue) Bytes(bo ...binary.ByteOrder) []byte {
	var a = make([]byte, 2)
	var b binary.ByteOrder = binary.BigEndian
	if len(bo) > 0 && bo[0] == binary.LittleEndian {
		b = binary.LittleEndian
	}
	b.PutUint16(a, uint16(t))
	return a
}

// TODO Add other validator
var (
	defaultCompression = defaultValidator{v: []TagValue{
		NoCompression, CCITTGroup31, CITTGroup3FaxT4, CITTGroup3FaxT6, LZW, JPEG, PackBits,
	}}
)

// Compression valid val
const (
	NoCompression TagValue = iota + 0x0001
	// CCITTGroup31 Hoffman
	CCITTGroup31
	CITTGroup3FaxT4
	CITTGroup3FaxT6
	LZW
	JPEG
	PackBits TagValue = 0x8005
)

func (t TagValue) Decode(data []byte) ([]byte, error) {

	switch t {
	case LZW:
		res := lzw.NewReader(bytes.NewReader(data), lzw.LSB, 8)
		defer res.Close()
		all, err := io.ReadAll(res)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return all, nil
	case JPEG:
		decode, err := jpeg.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return ImageByte(decode), nil

	default:
		return data, nil
	}
}

func ImageByte(img image.Image) []byte {
	s := img.Bounds().Size()
	x, y := s.X, s.Y
	var res = make([]byte, 0, x*y*3)
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			res = append(res, byte(r), byte(g), byte(b))
		}
	}
	return res
}

// ResolutionUnit  valid val
const (
	_ TagValue = iota
	NoUnit
	Inch
	Cm
)

// BitsPerSample
const (
	v1 TagValue = 1 << 2
	v2          = v1 << 1
)

// PhotoMetricInterpretation  valid val
const (
	_              = iota
	White TagValue = iota - 1
	Black
	RGB
	ColorIdx
	TransparentMixed
	CMYK
	YCbCr
	_
	Lab
	CFA TagValue = 0x8023
)

// ExtraSamples  valid val
const (
	e0 TagValue = iota
	Alpha
	NoAlpha
)

// FillOrder valid val
const (
	_ TagValue = iota
	LeftToRight
	RightToLeft
)

// GrayResponseUnit valid val
const (
	_ TagValue = iota
	OneTenth
	OneHun
	OneTho
	OneTenTho
	OneHunTho
)

// NextSubFileType valid val
const (
	_ TagValue = iota
	Full
	Low
	Multi
	Mixed
)

// Orient valid val
//
//	First Row and First column orientation
const (
	_ TagValue = iota
	UpL
	UpR
	BottomR
	BottomL
	LeftU
	RightU
	RightB
	LeftB
)

// PlanarConfiguration valid val
const (
	_ TagValue = iota
	SinglePlanar
	MultiPlanar
)

// SamplesPerPixel valid val
const (
	_ TagValue = iota
	BlackWhiteSample
	_
	RGBSample
)

// Threashholding valid val
const (
	_ TagValue = iota
	NoThreashholding
	// Dither also is halftone
	Dither
	ErrorDiffusion
)
