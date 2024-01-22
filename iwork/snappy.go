package iwork

import (
	"bytes"
	"github.com/klauspost/compress/snappy"
	"github.com/pkg/errors"
	"io"
)

func DecodeSnappy(reader io.Reader) (io.Reader, error) {

	var res []byte
	all, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	_, err = snappy.Decode(res, all)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return bytes.NewBuffer(res), nil
}
