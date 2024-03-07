package tiff

import (
	"encoding/binary"
	"log"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestName(t *testing.T) {
	open, err := os.Open("WTF.y.tif")
	if err != nil {
		panic(err)
	}

	tt, err := openTiff(open)
	if err != nil {
		log.Printf("%+v", err)
		panic("!!")
	}
	create2, err := os.Create("y.temp")
	if err != nil {
		panic(err)
	}

	bo := tt.ifh.bo
	spew.Dump(tt.ifh)
	for i := 0; i < len(tt.ifds); i++ {
		spew.Fprintf(create2, "------------------------------")
		dd := tt.ifds[i].de

		for i := range dd {
			//if Tag(tt.ifh.bo.Uint16(dd[i].tag[:])) == StripByteCounts {
			//	stripByteCounts = tt.ifh.bo.Uint32(dd[i].val[:])
			//}

			spew.Fprintf(create2, "{\n"+
				" tag: %s \n"+
				" val: %d %d %x  \n"+
				" offset: %x   \n"+
				"}\n", Tag(bo.Uint16(dd[i].tag[:])), bo.Uint16(dd[i].val[:]), bo.Uint32(dd[i].val[:]), dd[i].val[:], dd[i].offset[:])
			//fmt.Fprint(create, dd[i])

		}

	}
	log.Println(len(tt.vals))
	var aa = make([]uint16, 0, len(tt.vals))
	for _, val := range tt.vals {
		aa = append(aa, uint16(val))
	}

}

func Test01x(t *testing.T) {
	var a = []byte{00, 0x04, 00, 00}

	log.Println(binary.BigEndian.Uint32(a))
}

const a = 0x01f4

var xunit = []byte{00, 0x97, 0x33, 0x33, 00, 0x04, 00, 00}

func TestGenerate(t *testing.T) {

	bs, err := StringToRGBA("  HELLO,KNOW WEATHER!")
	if err != nil {
		log.Printf("%+v", err)
	}
	var dd []byte
	for i := range bs {
		if i%4 != 3 {
			dd = append(dd, bs[i])
		}
	}
	bs = dd
	var byteCount = make([]byte, 4)
	var firstIFDOffset = make([]byte, 4)
	var xuo, yuo, bp = make([]byte, 4), make([]byte, 4), make([]byte, 4)
	binary.BigEndian.PutUint32(byteCount, uint32(len(bs)))
	binary.BigEndian.PutUint32(firstIFDOffset, uint32(len(bs)+8))
	binary.BigEndian.PutUint32(xuo, uint32(len(bs)+14+17*12))
	binary.BigEndian.PutUint32(yuo, uint32(len(bs)+22+17*12))
	binary.BigEndian.PutUint32(bp, uint32(len(bs)+30+17*12))
	ti := tiff{
		ifh: ifh{
			bo:             binary.BigEndian,
			byteorder:      [2]byte{'M', 'M'},
			version:        [2]byte(DefaultTiff.Bytes()),
			firstIFDOffset: [4]byte(firstIFDOffset),
		},
		ifds: []*ifd{
			{num: [2]byte{00, 0x11},
				de: []*de{
					{
						tag:   [2]byte(ImageWidth.Bytes()),
						typ:   [2]byte(LONG.Bytes()),
						count: count01,
						val:   [4]byte{0, 0, 1, 0xf4},
					},
					{
						tag:   [2]byte(ImageLength.Bytes()),
						typ:   [2]byte(LONG.Bytes()),
						count: count01,
						val:   [4]byte{0, 0, 1, 0x2c},
					},
					{
						tag:   [2]byte(BitsPerSample.Bytes()),
						typ:   [2]byte(SHORT.Bytes()),
						count: count03,
						val:   [4]byte(bp),
					},
					{
						tag:   [2]byte(Compression.Bytes()),
						typ:   [2]byte(SHORT.Bytes()),
						count: count01,
						val:   short01,
					},
					{
						tag:   [2]byte(PhotoMetricInterpretation.Bytes()),
						typ:   [2]byte(SHORT.Bytes()),
						count: count01,
						val:   short02,
					},
					{
						tag:   [2]byte(FillOrder.Bytes()),
						typ:   [2]byte(SHORT.Bytes()),
						count: count01,
						val:   short01,
					},
					{
						tag:   [2]byte(StripOffsets.Bytes()),
						typ:   [2]byte(LONG.Bytes()),
						count: count01,
						val:   [4]byte{0, 0, 0, 8},
					},
					{
						tag:   [2]byte(Orient.Bytes()),
						typ:   [2]byte(SHORT.Bytes()),
						count: count01,
						val:   short01,
					},
					{
						tag:   [2]byte(SamplesPerPixel.Bytes()),
						typ:   [2]byte(SHORT.Bytes()),
						count: count01,
						val:   short03,
					},
					{
						tag:   [2]byte(RowsPerStrip.Bytes()),
						typ:   [2]byte(LONG.Bytes()),
						count: count01,
						val:   [4]byte{0, 0, 1, 0x2c},
					},
					{
						tag:   [2]byte(StripByteCounts.Bytes()),
						typ:   [2]byte(LONG.Bytes()),
						count: count01,
						val:   [4]byte(byteCount),
					},
					{
						tag:   [2]byte(XResolution.Bytes()),
						typ:   [2]byte(RATIONAL.Bytes()),
						count: count01,
						val:   [4]byte(xuo),
					},
					{
						tag:   [2]byte(YResolution.Bytes()),
						typ:   [2]byte(RATIONAL.Bytes()),
						count: count01,
						val:   [4]byte(yuo),
					},
					{
						tag:   [2]byte(PlanarConfiguration.Bytes()),
						typ:   [2]byte(SHORT.Bytes()),
						count: count01,
						val:   short01,
					},
					{
						tag:   [2]byte(ResolutionUnit.Bytes()),
						typ:   [2]byte(SHORT.Bytes()),
						count: count01,
						val:   short03,
					},
					{
						tag:   [2]byte(Artist.Bytes()),
						typ:   [2]byte(ASCII.Bytes()),
						count: count01,
						val:   [4]byte{'Z', 'E', 'R', 'O'},
					},
					{
						tag:   [2]byte(Copyright.Bytes()),
						typ:   [2]byte(ASCII.Bytes()),
						count: count01,
						val:   [4]byte{'@', 'K', 'W', 0},
					},
				},
			},
		},
		vals: bs,
	}
	bb := ti.byte()
	bb = append(bb, xunit...)
	bb = append(bb, xunit...)
	bb = append(bb, 0, 0x08, 0, 0x08, 0, 0x08)

	create, err := os.Create("WTF.tiff")
	if err != nil {
		panic(err)
	}
	_, err = create.Write(bb)
	if err != nil {
		panic(err)
	}

}

func TestTiff(t *testing.T) {
	bs, err := StringToRGBA("  HELLO,KNOW WEATHER!")
	if err != nil {
		log.Printf("%+v", err)
	}
	var dd []byte
	for i := range bs {
		if i%4 != 3 {
			dd = append(dd, bs[i])
		}
	}
	bs = dd
	ti := &Tiff{
		IFH: IFH{
			ByteOrder: MM,
			Version:   DefaultTiff,
		},
		Data: bs,
	}
	i := &IFD{}
	ti.addIFD(i)

	i.addDe(&DE{
		Tag:      Compression,
		DataType: SHORT,
		TagVal:   []Bytes{ComplementFour{NoCompression}},
	})
	i.addDe(&DE{
		Tag:      ImageWidth,
		DataType: LONG,
		TagVal:   []Bytes{Long(500)},
	})
	i.addDe(&DE{
		Tag:      ImageLength,
		DataType: LONG,
		TagVal:   []Bytes{Long(300)},
	})
	i.addDe(&DE{
		Tag:      Copyright,
		DataType: ASCII,
		count:    11,
		TagVal:   []Bytes{String("@KW")},
	})
	i.addDe(&DE{
		Tag:      ResolutionUnit,
		DataType: SHORT,
		TagVal:   []Bytes{ComplementFour{Cm}},
	})
	i.addDe(&DE{
		Tag:      BitsPerSample,
		DataType: SHORT,
		TagVal:   []Bytes{Short(8), Short(8), Short(8)},
	})
	i.addDe(&DE{
		Tag:      PhotoMetricInterpretation,
		DataType: SHORT,
		TagVal:   []Bytes{ComplementFour{RGB}},
	})
	i.addDe(&DE{
		Tag:      FillOrder,
		DataType: SHORT,
		TagVal:   []Bytes{ComplementFour{LeftToRight}},
	})
	i.addDe(&DE{
		Tag:      Artist,
		DataType: ASCII,
		count:    4,
		TagVal:   []Bytes{String("ZERO")},
	})
	i.addDe(&DE{
		Tag:      XResolution,
		DataType: RATIONAL,
		TagVal:   []Bytes{Rational{9909043, 262144}},
	})
	i.addDe(&DE{
		Tag:      YResolution,
		DataType: RATIONAL,
		TagVal:   []Bytes{Rational{9909043, 262144}},
	})
	i.addDe(&DE{
		Tag:      Orient,
		DataType: SHORT,
		TagVal:   []Bytes{ComplementFour{UpL}},
	})
	i.addDe(&DE{
		Tag:      SamplesPerPixel,
		DataType: SHORT,
		TagVal:   []Bytes{ComplementFour{RGBSample}},
	})
	i.addDe(&DE{
		Tag:      RowsPerStrip,
		DataType: LONG,
		TagVal:   []Bytes{Long(300)},
	})
	i.addDe(&DE{
		Tag:      StripByteCounts,
		DataType: LONG,
		TagVal:   []Bytes{Long(len(bs))},
	})
	i.addDe(&DE{
		Tag:      PlanarConfiguration,
		DataType: SHORT,
		TagVal:   []Bytes{ComplementFour{SinglePlanar}},
	})

	i.addDe(&DE{
		Tag:      StripOffsets,
		DataType: LONG,
		TagVal:   []Bytes{Long(8)},
	})

	create, err := os.Create("WTF.y.tif")
	if err != nil {
		panic(err)
	}

	err = ti.Save(create)
	if err != nil {
		panic(err)
	}
}

func TestSimpleTiff_Save(t *testing.T) {
	bs, err := StringToRGBA("  HELLO!")
	if err != nil {
		log.Printf("%+v", err)
	}
	var dd []byte
	for i := range bs {
		if i%4 != 3 {
			dd = append(dd, bs[i])
		}
	}
	bs = dd

	err = CreateSimpleTiff("Hello", bs, 500, 300).Save()
	if err != nil {
		panic(err)
	}
}
