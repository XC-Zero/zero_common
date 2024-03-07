package tiff

import (
	"encoding/binary"
	"io"
	"log"

	"github.com/pkg/errors"
)

// Deprecated
type tiff struct {
	ifh         ifh
	ifds        []*ifd
	valueOffset [4]byte
	vals        []byte
	aux         axuSpace
}

// Deprecated
type ifh struct {
	bo             binary.ByteOrder
	byteorder      [2]byte
	version        [2]byte
	firstIFDOffset [4]byte
}

// Deprecated
type ifd struct {
	num           [2]byte
	de            []*de
	nextIFDOffset [4]byte
}

// Deprecated
type de struct {
	tag    [2]byte
	typ    [2]byte
	count  [4]byte
	val    [4]byte
	offset [4]byte
}

// Deprecated
type axuSpace struct {
	data []byte
}

var (
	ReadIFHErr = errors.New("Read file header failed!")
	BrErr      = errors.New("Header byte order is incorrect!")
)

// Deprecated
func openTiff(r io.ReadSeeker) (*tiff, error) {
	var bs = make([]byte, 8)
	_, err := r.Seek(0, 0)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	all, err := r.Read(bs)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if all != 8 {
		return nil, errors.WithStack(ReadIFHErr)
	}
	br := bs[:2]
	var bo binary.ByteOrder
	switch ByteOrder(br) {
	case II:
		bo = binary.LittleEndian
	case MM:
		bo = binary.BigEndian
	default:
		return nil, errors.WithStack(BrErr)
	}

	t := &tiff{
		ifh: ifh{
			bo:             bo,
			byteorder:      [2]byte(br),
			version:        [2]byte(bs[2:4]),
			firstIFDOffset: [4]byte(bs[4:8]),
		},
	}
	nxt, err := t.readIFD(r, int64(bo.Uint32(bs[4:8])))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var i = 0
	for nxt != 0 {
		log.Printf("nxt:%x", nxt)
		nxt, err = t.readIFD(r, int64(nxt))
		if err != nil {
			return nil, errors.WithStack(err)
		}
		i++
		if i > 5 {
			break
		}
	}
	err = t.readData(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return t, nil
}

// Deprecated
func (t *tiff) readIFD(r io.ReadSeeker, offset int64) (uint32, error) {
	_, err := r.Seek(offset, 0)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	var n = make([]byte, 2)
	_, err = r.Read(n)
	if err != nil && err != io.EOF {
		return 0, errors.WithStack(err)
	}
	num := int(t.ifh.bo.Uint16(n))
	var bs = make([]byte, num*12+4)
	_, err = r.Read(bs)
	if err != nil && err != io.EOF {
		return 0, errors.WithStack(err)
	}
	var _ifd = &ifd{
		num: [2]byte(n),
	}
	log.Println(num)
	log.Printf("%d %x", num, num*12+4)
	for i := 0; i < num; i++ {
		start := i * 12
		var d = de{
			tag:   [2]byte(bs[start : start+2]),
			typ:   [2]byte(bs[start+2 : start+4]),
			count: [4]byte(bs[start+4 : start+8]),
		}
		vo := [4]byte(bs[start+8 : start+12])
		if t.ifh.bo.Uint16(d.tag[:]) == uint16(StripOffsets) {
			t.valueOffset = vo
		}
		count := t.ifh.bo.Uint32(d.count[:])
		count *= DataType(t.ifh.bo.Uint16(d.typ[:])).Len()
		if count < 5 {
			d.val = vo
		} else {
			d.offset = vo
		}
		_ifd.de = append(_ifd.de, &d)
	}
	bs = make([]byte, 4)
	_, err = r.Read(bs)
	if err != nil && err != io.EOF {
		return 0, errors.WithStack(err)
	}
	_ifd.nextIFDOffset = [4]byte(bs)
	t.ifds = append(t.ifds, _ifd)
	nextOffset := t.ifh.bo.Uint32(bs)

	if nextOffset != 0 && nextOffset != 65535 {
		return nextOffset, nil
	}

	return 0, nil
}

// Deprecated
func (t *tiff) readData(r io.ReadSeeker) error {
	_, err := r.Seek(int64(t.ifh.bo.Uint16(t.valueOffset[:])), 0)
	if err != nil {
		return errors.WithStack(err)
	}
	t.vals, err = io.ReadAll(r)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Deprecated
// FIXME Deal with oversize data offset, there are just normal data...
func (t *tiff) save(w io.Writer) error {

	_, err := w.Write(t.byte())
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Deprecated
func (t *tiff) byte() []byte {
	var bs []byte

	bs = append(bs, t.ifh.byteorder[:]...)
	bs = append(bs, t.ifh.version[:]...)
	bs = append(bs, t.ifh.firstIFDOffset[:]...)
	bs = append(bs, t.vals...)

	for i := range t.ifds {
		bs = append(bs, t.ifds[i].num[:]...)
		for j := range t.ifds[i].de {
			de := t.ifds[i].de[j]
			bs = append(bs, de.tag[:]...)
			bs = append(bs, de.typ[:]...)
			bs = append(bs, de.count[:]...)
			bs = append(bs, de.val[:]...)
		}
		bs = append(bs, t.ifds[i].nextIFDOffset[:]...)

	}
	return bs
}

//func (t *tiff) ToTiff() {
//	var tagMap = make(map[Tag]any)
//	bo := t.ifh.bo
//	for i := range t.ifds[0].de {
//		d := t.ifds[0].de[i]
//		dt := DataType(bo.Uint32(d.typ[:]))
//		var v any
//		switch dt {
//		case SHORT:
//			v = bo.Uint16(d.val[:])
//		case LONG:
//			v = bo.Uint32(d.val[:])
//		case BYTE:
//			v = string(d.val[:])
//		default:
//			v = nil
//		}
//		tagMap[Tag(bo.Uint16(d.tag[:]))] = v
//	}
//
//	Tiff{
//		ImageWidth:     tagMap[ImageWidth].(uint32),
//		ImageLength:    tagMap[ImageLength].(uint32),
//		Compression:    TagValue(tagMap[Compression].(uint16)),
//		ResolutionUnit: 0,
//		software:       "",
//		datetime:       time.Time{},
//		xResolution:    0,
//		yResolution:    0,
//		stripOffset:    0,
//		stripByteCount: 0,
//		colorMap:       nil,
//	}
//}
