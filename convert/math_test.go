package convert

import (
	"log"
	"testing"
)

func TestTrimRightZeroRetain(t *testing.T) {
	log.Println(TrimRightZeroRetain("26.0583", 2))
}
