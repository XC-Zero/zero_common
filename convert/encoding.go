package convert

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"strings"
)

type Coding int

const (
	GB18030 Coding = iota + 2
	GBK
	GB2312
	HZGB2312
	EUCCN
	BIG5
	BIGFIVE
	UTF8
	ISO8859_1
	ISO8859_2
	ISO2022JP
	CP932
	CP936
	ShiftJIS
	Windows31J
	Undefined
)

func (c Coding) String() string {
	switch c {
	case GB2312:
		return "GB2312"
	case ISO8859_1:
		return "ISO8859_1"
	case ISO8859_2:
		return "ISO8859_2"
	case GBK:
		return "GBK"
	case CP936:
		return "CP936"
	case GB18030:
		return "GB18030"
	case UTF8:
		return "UTF8"
	case ShiftJIS:
		return "ShiftJIS"
	case CP932:
		return "CP932"
	case Windows31J:
		return "Windows31J"
	case ISO2022JP:
		return "ISO2022JP"
	default:
		return "Undefined"
	}
}

func GetEncodingByName(name string) Coding {
	name = strings.ToUpper(name)
	switch name {
	case "GB2312":
		return GB2312
	case "ISO8859_1":
		return ISO8859_1
	case "ISO8859_2":
		return ISO8859_2
	case "GBK":
		return GBK
	case "CP936":
		return CP936
	case "GB18030":
		return GB18030
	case "UTF8":
		return UTF8
	case "ShiftJIS":
		return ShiftJIS
	case "CP932":
		return CP932
	case "Windows31J":
		return Windows31J
	case "ISO2022JP":
		return ISO2022JP
	default:
		return Undefined
	}
}

func (c Coding) Coding() encoding.Encoding {
	return GetEncoding(c)
}

func GetEncoding(coding Coding) encoding.Encoding {
	switch coding {
	case GB2312, HZGB2312, EUCCN:
		return simplifiedchinese.HZGB2312
	case GBK, CP936:
		return simplifiedchinese.GBK
	case GB18030:
		return simplifiedchinese.GB18030
	case BIG5, BIGFIVE:
		return traditionalchinese.Big5
	case ISO8859_1:
		return charmap.ISO8859_1
	case ISO8859_2:
		return charmap.ISO8859_2
	case CP932, ShiftJIS, Windows31J:
		return japanese.ShiftJIS
	case ISO2022JP:
		return japanese.ISO2022JP
	case UTF8:
		return encoding.Nop
	}
	return encoding.Nop
}

func Transform(input io.Reader, inputCoding Coding, output io.Writer, outputCoding Coding) {
	from := transform.NewReader(input, inputCoding.Coding().NewDecoder())
	to := transform.NewWriter(output, outputCoding.Coding().NewDecoder())
	_, err := io.Copy(to, from)
	if err != nil {
		log.Println(err)
	}
}

func GB2312ToUTF8(input io.Reader, output io.Writer) {
	from := transform.NewReader(input, GB2312.Coding().NewDecoder())
	to := transform.NewWriter(output, UTF8.Coding().NewDecoder())
	_, err := io.Copy(to, from)
	if err != nil {
		log.Println(err)
	}
}

func GB18030ToUTF8(input io.Reader, output io.Writer) {
	from := transform.NewReader(input, GB18030.Coding().NewDecoder())
	to := transform.NewWriter(output, UTF8.Coding().NewDecoder())
	_, err := io.Copy(to, from)
	if err != nil {
		log.Println(err)
	}
}
