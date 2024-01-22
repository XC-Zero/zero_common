package iwork

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jhump/protoreflect/desc/protoparse"
)

func DecodeProtobuf(filePath string) {
	Parser := protoparse.Parser{}
	descs, err := Parser.ParseFiles(filePath)
	if err != nil {
		fmt.Printf("ParseFiles err=%v", err)
		return
	}
	for i := 0; i < len(descs); i++ {
		spew.Dump(descs[i].String())
	}
}
