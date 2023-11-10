package monitor

import (
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestGetAllInfo(t *testing.T) {
	spew.Dump(GetAllInfo())
}
