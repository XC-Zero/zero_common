package convert

import (
	"io"
	"os"
	"testing"
)

func TestToZip(t *testing.T) {
	open, err := os.Open("A:\\download\\diiv_4月.xlsx")
	if err != nil {
		panic(err)
	}

	open2, err := os.Open("A:\\download\\diiv_4.xlsx")
	if err != nil {
		panic(err)
	}
	a, err := os.Create("test.zip")
	if err != nil {
		panic(err)
	}
	b, err := os.Create("test2.zip")
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(a, FileToZip(open2, open))
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(b, ReaderToZip(map[string]io.Reader{"test01x.xlsx": open2, "./te/test02x.xlsx": open}))
	if err != nil {
		panic(err)
	}
}
