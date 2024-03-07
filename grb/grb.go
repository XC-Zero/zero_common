package grb

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/nilsmagnus/grib/griblib"
	"github.com/pkg/errors"
	"os"
)

type Grb struct {
}

func NewGrb(filepath string) (*Grb, error) {
	open, err := os.Open(filepath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	messages, err := griblib.ReadNMessages(open)

	for i := range messages {
		msg := messages[i]
		spew.Dump(msg)
	}
	return nil, nil
}
