package convert

import (
	"log"
	"testing"
)

func TestStrSliceToIntSlice(t *testing.T) {
}

func TestStringToIntSlice(t *testing.T) {
	log.Println(StringToIntSlice("9415", ","))
}
