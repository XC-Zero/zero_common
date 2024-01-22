package iwork

import (
	"github.com/davecgh/go-spew/spew"
	"io"
	"os"
	"testing"
)

func TestDecodeProtobuf(t *testing.T) {
	filePath := "C:\\Users\\XC\\Desktop\\Metadata.iwa"
	open, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	//bb := []byte{0x82, 0x4b, 0x4e, 0x50, 0x53, 0x4e, 0x41, 0x50, 0x50, 0x59}
	all, err := io.ReadAll(open)
	if err != nil {
		panic(err)
	}
	spew.Dump(string(all))
	//bb = append(bb, all...)
	//create, err := os.Create(filePath + "xxx")
	//if err != nil {
	//	panic(err)
	//}
	//_, err = io.Copy(create, bytes.NewBuffer(bb))
	//if err != nil {
	//	panic(err)
	//}
	//snappy, err := DecodeSnappy(bytes.NewBuffer(bb))
	//if err != nil {
	//	panic(err)
	//}
	//currentDir, err := os.Getwd()
	//if err != nil {
	//	panic(err)
	//
	//}
	//filePath = path.Join(currentDir, "temp01x")
	//create, err = os.Create(filePath)
	//if err != nil {
	//	panic(err)
	//}
	//_, err = io.Copy(create, snappy)
	//if err != nil {
	//	panic(err)
	//}
	DecodeProtobuf(filePath)

}
