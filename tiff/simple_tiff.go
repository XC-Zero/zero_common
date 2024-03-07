package tiff

import (
	"os"

	"github.com/pkg/errors"
)

type SimpleTiff struct {
	name string
	t    *Tiff
}

// CreateSimpleTiff imageData only support RGB max size is 2>>32 ,seem like 4G
func CreateSimpleTiff(name string, imageData []byte, width, length int) *SimpleTiff {
	s := &SimpleTiff{
		name: name,
		t: &Tiff{
			IFH: IFH{
				ByteOrder: MM,
				Version:   DefaultTiff,
			},
			IFDs: nil,
			Data: imageData,
		},
	}
	i := &IFD{}
	s.t.addIFD(i)
	i.addDe(&DE{
		Tag:      Compression,
		DataType: SHORT,
		TagVal:   []Bytes{ComplementFour{NoCompression}},
	})
	i.addDe(&DE{
		Tag:      ImageWidth,
		DataType: LONG,
		TagVal:   []Bytes{Long(width)},
	})
	i.addDe(&DE{
		Tag:      ImageLength,
		DataType: LONG,
		TagVal:   []Bytes{Long(length)},
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
		TagVal:   []Bytes{Long(length)},
	})
	i.addDe(&DE{
		Tag:      StripByteCounts,
		DataType: LONG,
		TagVal:   []Bytes{Long(len(imageData))},
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

	return s

}

// PutAuxTag Will cover same name tag
func (s *SimpleTiff) PutAuxTag(tag Tag, dataType DataType, val []Bytes) *SimpleTiff {
	i := s.t.IFDs[0]
	i.addDe(&DE{
		Tag:      tag,
		DataType: dataType,
		TagVal:   val,
	})
	return s
}

func (s *SimpleTiff) Save() error {
	create, err := os.Create(s.name + ".tif")
	if err != nil {
		return errors.WithStack(err)
	}
	err = s.t.Save(create)
	if err != nil {
		return err
	}
	return nil
}
