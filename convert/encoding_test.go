package convert

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
)

func TestTransform(t *testing.T) {
	open, err := os.Open("A:\\xccjm\\Documents\\a.txt")
	if err != nil {
		panic(err)

	}
	tt, err := os.Create("a_utf8_2.txt")
	if err != nil {
		panic(err)
	}
	defer tt.Close()
	a := bytes.NewBuffer(nil)
	GB18030ToUTF8(open, a)
	_, m := ToJson(a)
	marshal, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewBuffer(marshal)
	_, err = io.Copy(tt, buffer)
	if err != nil {
		panic(err)
	}

}
