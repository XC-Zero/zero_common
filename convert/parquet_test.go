package convert

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestTxtToParquet(t *testing.T) {
	file, err := os.Open("A:\\xccjm\\Documents\\a.txt")
	if err != nil {
		panic(err)
	}

	buff := bytes.NewBuffer(nil)
	GB18030ToUTF8(file, buff)
	name := strings.Split(file.Name(), "\\")
	stat, err := file.Stat()
	if err != nil {
		return
	}
	stat.Name()
	col := TxtToParquet(buff, name[len(name)-1], 1)
	if err != nil {
		panic(err)
	}
	log.Println("colcc is :", col)

	//xx, err := os.Create("xxx.parquet")
	//if err != nil {
	//	panic(err)
	//}
	//_, err = io.Copy(xx, f)
	//if err != nil {
	//	panic(err)
	//}
}
