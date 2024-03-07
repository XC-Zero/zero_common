package tiff

import (
	"encoding/binary"
	"io"
	"log"
	"slices"
	"time"

	"github.com/pkg/errors"
)

type Tiff struct {
	IFH       IFH
	IFDs      []*IFD
	Data      []byte
	auxSpace  []byte
	auxOffset uint32
}

func (t *Tiff) addIFD(ifd *IFD) {
	ifd.tiff = t
	ifd.len = 6
	t.IFDs = append(t.IFDs, ifd)
}

func (t *Tiff) Bytes(_ ...binary.ByteOrder) []byte {
	var bs []byte
	var now = time.Now()
	t.IFH.tiff = t
	bs = append(bs, t.IFH.Bytes()...)
	t.auxOffset += uint32(len(t.Data))

	t.auxOffset += uint32(len(bs))

	for i := 0; i < len(t.IFDs); i++ {
		t.auxOffset += t.IFDs[i].num*12 + 6
	}
	bs = append(bs, t.Data...)

	log.Printf("Generate IFH %s", time.Since(now))
	bo := t.IFH.ByteOrder.Bo()
	for i := range t.IFDs {
		t.IFDs[i].idx = uint32(i)
		t.IFDs[i].startOffset = uint32(len(bs))
		bs = append(bs, t.IFDs[i].Bytes(bo)...)
	}
	log.Printf("Generate IFD %s", time.Since(now))

	bs = append(bs, t.auxSpace...)
	return bs
}

type IFH struct {
	tiff      *Tiff
	ByteOrder ByteOrder
	Version   Version
}

func (i *IFH) Bytes(_ ...binary.ByteOrder) []byte {
	var bs = make([]byte, 0, 8)
	bo := i.ByteOrder.Bo()
	bs = append(bs, i.ByteOrder.Bytes()...)
	bs = append(bs, i.Version.Bytes(bo)...)
	var offset = make([]byte, 4)
	bo.PutUint32(offset, uint32(8+len(i.tiff.Data)))
	bs = append(bs, offset...)
	return bs

}

type IFD struct {
	tiff        *Tiff
	idx         uint32
	startOffset uint32
	len         int
	num         uint32
	DEs         []*DE
}

func (i *IFD) addDe(de *DE) {
	if de.count == 0 {
		de.count = Long(len(de.TagVal))
	}
	de.tiff = i.tiff
	for j, e := range i.DEs {
		if e.Tag == de.Tag {
			i.DEs[j].TagVal = de.TagVal
			i.DEs[j].DataType = de.DataType
			i.DEs[j].count = de.count
			return
		}
	}
	i.len += 12

	i.DEs = append(i.DEs, de)
	i.num++
	slices.SortFunc(i.DEs, DESort)

}

func (i *IFD) Len() int {
	return i.len
}

func (i *IFD) Bytes(bo ...binary.ByteOrder) []byte {
	var l = 2 + i.num*12 + 4
	var bs = make([]byte, 2, l)
	var b binary.ByteOrder = binary.BigEndian
	if len(bo) > 0 && bo[0] == binary.LittleEndian {
		b = binary.LittleEndian
	}
	b.PutUint16(bs, uint16(i.num))
	for j := 0; j < int(i.num); j++ {
		bs = append(bs, i.DEs[j].Bytes(b)...)
	}
	var nxt = make([]byte, 4)
	if len(i.tiff.IFDs) == int(i.idx+1) {
		b.PutUint32(nxt, 0)
	} else {
		b.PutUint32(nxt, i.startOffset+l)
	}
	bs = append(bs, nxt...)
	return bs
}

type DE struct {
	tiff     *Tiff
	Tag      Tag
	DataType DataType
	count    Long
	needAux  bool
	TagVal   []Bytes
}

func (d *DE) NeedAux() {
	if d.DataType.Len()*uint32(len(d.TagVal)) > 4 {
		d.needAux = true
	}
}

func (d *DE) Bytes(bo ...binary.ByteOrder) []byte {
	d.NeedAux()
	var b binary.ByteOrder = binary.BigEndian
	if len(bo) > 0 && bo[0] == binary.LittleEndian {
		b = binary.LittleEndian
	}
	var bs = make([]byte, 0, 12)
	bs = append(bs, d.Tag.Bytes(b)...)
	bs = append(bs, d.DataType.Bytes(b)...)
	bs = append(bs, d.count.Bytes(b)...)

	var val []byte
	for i := range d.TagVal {
		val = append(val, d.TagVal[i].Bytes(b)...)
	}
	var tagVal = make([]byte, 4)
	if d.needAux {
		var auxOffset = len(d.tiff.auxSpace)
		b.PutUint32(tagVal, uint32(auxOffset)+d.tiff.auxOffset)
		d.tiff.auxSpace = append(d.tiff.auxSpace, val...)
	} else {
		for i := range val {
			if i >= 4 {
				break
			}
			tagVal[i] = val[i]
		}
	}
	bs = append(bs, tagVal...)
	return bs
}

func DESort(a, b *DE) int {
	if a.Tag < b.Tag {
		return -1
	}
	if a.Tag > b.Tag {
		return 1
	}
	return 0
}

func (t *Tiff) Save(w io.Writer) error {
	_, err := w.Write(t.Bytes())
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

var (
	short01 = [4]byte{0, 1, 0, 0}
	short02 = [4]byte{0, 2, 0, 0}
	short03 = [4]byte{0, 3, 0, 0}
	short04 = [4]byte{0, 4, 0, 0}
	short08 = [4]byte{0, 8, 0, 0}

	count01 = [4]byte{0, 0, 0, 1}
	count02 = [4]byte{0, 0, 0, 2}
	count03 = [4]byte{0, 0, 0, 3}
)

//type Tiff struct {
//	ImageWidth     Long
//	ImageLength    Long
//	Compression    TagValue
//	ResolutionUnit TagValue
//	Software       string
//	Datetime       time.Time
//	XResolution    Rational
//	YResolution    Rational
//	StripOffset    Long
//	StripByteCount Long
//	ExtraTag       map[Tag][4]byte
//	Data           []rgb
//}
//
//type rgb struct {
//	r, g, b Byte
//}
//
////func Open(path string) *Tiff {
////	binary.Read()
////}
////
////func New() *Tiff {
////
////}
////
////func Load[T float64 | float32 | int | int8 | int16 | int32 | int64](data [][]T) {
////
////}
//
//// Save TODO Save
//func (t *Tiff) Save(path string) error {
//	var (
//		byteCount = make([]byte, 4)
//		//width     = make([]byte, 4)
//		length = make([]byte, 4)
//	)
//
//	binary.BigEndian.PutUint32(byteCount, uint32(len(t.Data)*3))
//	ti := tiff{
//		ifh: ifh{
//			bo:             binary.BigEndian,
//			byteorder:      [2]byte{'M', 'M'},
//			version:        DefaultTiff.Bytes(),
//			firstIFDOffset: [4]byte(byteCount),
//		},
//		ifds: []*ifd{
//			{num: [2]byte{00, 0x0f},
//				de: []*de{
//					{
//						tag:   ImageWidth.Bytes(),
//						typ:   SHORT.Bytes(),
//						count: count01,
//						val:   t.ImageWidth.Bytes(),
//					},
//					{
//						tag:   ImageLength.Bytes(),
//						typ:   SHORT.Bytes(),
//						count: count01,
//						val:   [4]byte(length),
//					},
//					{
//						tag:   BitsPerSample.Bytes(),
//						typ:   SHORT.Bytes(),
//						count: count01,
//						val:   count03,
//					},
//					{
//						tag:   Compression.Bytes(),
//						typ:   SHORT.Bytes(),
//						count: count01,
//						val:   NoCompression.Bytes(),
//					},
//					{
//						tag:   PhotoMetricInterpretation.Bytes(),
//						typ:   SHORT.Bytes(),
//						count: count01,
//						val:   RGB.Bytes(),
//					},
//					{
//						tag:   FillOrder.Bytes(),
//						typ:   SHORT.Bytes(),
//						count: count01,
//						val:   LeftToRight.Bytes(),
//					},
//					{
//						tag:   StripOffsets.Bytes(),
//						typ:   LONG.Bytes(),
//						count: count01,
//						val:   [4]byte{0, 0, 0, 2},
//					},
//					{
//						tag:   Orient.Bytes(),
//						typ:   LONG.Bytes(),
//						count: count01,
//						val:   UpL.Bytes(),
//					},
//					{
//						tag:   SamplesPerPixel.Bytes(),
//						typ:   SHORT.Bytes(),
//						count: count01,
//						val:   count03,
//					},
//					{
//						tag:   RowsPerStrip.Bytes(),
//						typ:   LONG.Bytes(),
//						count: count01,
//						val:   count03,
//					},
//					{
//						tag:   SamplesPerPixel.Bytes(),
//						typ:   SHORT.Bytes(),
//						count: count01,
//						val:   count03,
//					},
//					{
//						tag:   SamplesPerPixel.Bytes(),
//						typ:   SHORT.Bytes(),
//						count: count01,
//						val:   count03,
//					},
//					{
//						tag:   SamplesPerPixel.Bytes(),
//						typ:   SHORT.Bytes(),
//						count: count01,
//						val:   count03,
//					},
//				},
//			},
//		},
//		valueOffset: [4]byte{},
//		vals:        nil,
//	}
//	create, err := os.Create(path)
//	if err != nil {
//		return errors.WithStack(err)
//	}
//	err = ti.save(create)
//	if err != nil {
//		return errors.WithStack(err)
//	}
//	return nil
//}
