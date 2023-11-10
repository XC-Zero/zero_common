package convert

import (
	"log"
	"strings"
	"testing"
)

func TestDealX(t *testing.T) {
	var a = "正常"
	var bb []byte

	for b := range a {
		bb = append(bb, byte(b))
	}
	bb = append(bb, '\xE3')
	bb = append(bb, '\xA1')

	log.Println(strings.ReplaceAll(string(DealX(bb)), "¡ã", "??"))
}
