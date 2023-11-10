package convert

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"os"
)

func FileToZip(files ...*os.File) *bytes.Buffer {
	res := bytes.NewBuffer([]byte{})
	zw := zip.NewWriter(res)
	defer zw.Close()
	for i := 0; i < len(files); i++ {
		var now = files[i]
		stat, err := now.Stat()
		if err != nil {
			panic(err)
		}
		name := stat.Name()
		log.Println(name)
		entry, err := zw.Create(name)
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(entry, now)
		if err != nil {
			panic(err)
		}

	}

	return res
}

func ReaderToZip(files map[string]io.Reader) *bytes.Buffer {
	res := bytes.NewBuffer([]byte{})
	zw := zip.NewWriter(res)
	defer zw.Close()

	for name, file := range files {
		entry, err := zw.Create(name)
		if err != nil {
			continue
		}
		io.Copy(entry, file)
	}

	return res
}
